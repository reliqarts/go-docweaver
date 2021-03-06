package docweaver

import (
	"fmt"
	cp "github.com/otiai10/copy"
	"os"
	"os/exec"
	"strings"
)

type DocHandler interface {
	GetDocsDir() string
}

type Publisher interface {
	Cleaner
	DocHandler
	Publish(productKey string, source string, shouldUpdate bool)
	// PublishFromSources publishes all documentation configured in sources file. i.e. env: DW_SOURCES_FILE
	PublishFromSources() (int, error)
}

type Updater interface {
	DocHandler
	Update(productKeys ...string)
	UpdateAll()
}

// UpdaterPublisher is a hybrid publisher/publisher.
type UpdaterPublisher interface {
	Updater
	Publisher
}

type publisher struct {
	repo ProductRepository
}

var mainVersions = []string{versionMaster, versionMain}

// GetPublisher returns the default instance of UpdaterPublisher.
func GetPublisher() UpdaterPublisher {
	return &publisher{
		repo: &productRepository{getDocsDir()},
	}
}

// GetPublisherWithDocsDir returns an instance of UpdaterPublisher with the provided [dir].
func GetPublisherWithDocsDir(docsDir string) UpdaterPublisher {
	return &publisher{repo: &productRepository{docsDir}}
}

func (p *publisher) Publish(productKey string, source string, shouldUpdate bool) {
	err := p.publish(productRoot{ParentDir: p.repo.GetDir(), Key: productKey, Source: source}, shouldUpdate)
	if err == nil {
		log(lInfo, "Successfully published product: `%s`.", productKey)
	}
}

func (p *publisher) GetDocsDir() string {
	return p.repo.GetDir()
}

func (p *publisher) PublishFromSources() (int, error) {
	sc, err := readSources()
	if err != nil {
		return -1, simpleError{fmt.Sprintf("Failed to publish documents from sources file. %s", err)}
	}

	for _, s := range sc.Sources {
		p.Publish(s.Key, s.Url, true)
	}

	return len(sc.Sources), nil
}

func (p *publisher) publish(pr productRoot, shouldUpdate bool) error {
	baseVersion := ""
	log(lInfo, "Publishing product: `%s`\n", pr.Key)
	log(lInfo, "Product root: %s\n", pr)

	for _, mv := range mainVersions {
		if err := p.publishProductVersion(pr, mv, true); err == nil {
			baseVersion = mv
			break
		}
	}

	if baseVersion == "" {
		return p.getBVMErr(pr.Key)
	}

	tags, err := p.listProductTags(pr, baseVersion)
	if err != nil {
		log(lError, "Failed to list tags. %s\n", err)
		return err
	}

	for _, tag := range tags {
		if err := p.publishProductVersion(pr, tag, shouldUpdate); err != nil {
			log(lWarn, "Failed to publish/update Tag `%s`.", tag)
		}
	}

	return nil
}

func (p *publisher) publishProductVersion(pr productRoot, version string, update bool) error {
	prFullPath := pr.filePath()
	versionTemp := versionTempName(version)
	verPath := pr.versionFilePath(version)
	verPathTemp := pr.versionFilePath(versionTemp)

	if _, err := os.Stat(prFullPath); os.IsNotExist(err) {
		if err = os.MkdirAll(prFullPath, 0755); err != nil {
			return err
		}
	}

	if _, err := os.Stat(verPath); !os.IsNotExist(err) {
		if !update {
			log(
				lInfo,
				"Version `%s` already exists for product `%s`. Update not requested. Skipped.\n",
				version,
				pr.Key,
			)
			return nil
		}
	}

	_ = removeDir(verPathTemp)
	log(lInfo, "Executing version clone `%s` into: `%s`\n", version, prFullPath)
	cmd := exec.Command("git", "clone", "--branch", version, pr.Source, versionTemp)
	cmd.Dir = pr.filePath()
	if _, err := cmd.Output(); err != nil {
		log(lError, "Failed to execute version clone `%s` into: `%s`. %s\n", version, prFullPath, err)
		return err
	}

	_ = removeDir(verPath)
	if err := os.Rename(verPathTemp, verPath); err != nil {
		log(lError, "Failed to rename version from temp file `%s` to `%s` after cloning update. %s\n", verPathTemp, verPath, err)
		return err
	}

	if err := p.publishVersionAssets(pr, version); err != nil {
		log(lError, "Failed to publish assets for version `%s`. %s\n", version, err)
	}

	return nil
}

func (p *publisher) listProductTags(pr productRoot, baseVersion string) ([]string, error) {
	var tags []string
	log(lInfo, "Listing tags for product `%s` using base version `%s`.\n", pr.Key, baseVersion)
	cmd := exec.Command("git", "tag", "-l")
	cmd.Dir = pr.versionFilePath(baseVersion)
	out, err := cmd.Output()
	if err != nil {
		log(lError, "Failed to list tags for product `%s` using base version `%s`.\n", pr.Key, baseVersion)
		return nil, err
	}

	for _, t := range strings.Split(string(out), "\n") {
		if t != "" {
			tags = append(tags, t)
		}
	}

	return tags, nil
}

func (p *publisher) Update(productKeys ...string) {
	for _, productName := range productKeys {
		log(lInfo, "Updating product: `%s`\n", productName)
		pr := productRoot{ParentDir: p.repo.GetDir(), Key: productName}
		baseVersion := ""
		source := ""

		for _, mv := range mainVersions {
			if _, err := os.Stat(pr.versionFilePath(mv)); !os.IsNotExist(err) {
				baseVersion = mv
				break
			}
		}

		if baseVersion == "" {
			log(lError, p.getBVMErr(productName).Error())
			continue
		}

		cmd := exec.Command("bash", "-c", "git remote show origin | grep Fetch")
		cmd.Dir = pr.versionFilePath(baseVersion)
		out, err := cmd.Output()
		if err != nil {
			log(lError, "Failed to determine fetch URL of origin for product `%s` using base version `%s`. %s\n", pr.Key, baseVersion, err)
			return
		}
		outSplit := strings.Split(string(out), ": ")
		if len(outSplit) == 2 {
			source = strings.TrimSpace(outSplit[1])
		}
		if source == "" {
			log(lError, "Could not determine source for product `%s`.\n", pr.Key)
			return
		}
		pr.Source = source

		err = p.publish(pr, true)
		if err == nil {
			log(lInfo, "Successfully updated product: `%s`.", productName)
		}
	}
}

func (p *publisher) UpdateAll() {
	productNames, err := p.repo.ListProductKeys()
	if err != nil {
		log(lError, err.Error())
	}
	if len(productNames) == 0 {
		log(lInfo, "No products found to update.")
		return
	}

	log(lInfo, "Updating the following products:", productNames)
	p.Update(productNames...)
}

// CleanTempVersions removes all temporary documentation versions. Only returns the last error that occurred.
func (p *publisher) CleanTempVersions() (lastErr error) {
	return p.repo.CleanTempVersions()
}

// getBVMErr generates a base version missing error with provided productName.
func (p *publisher) getBVMErr(productName string) error {
	return simpleError{fmt.Sprintf(
		"Base version for product %s could not be determined. Was not found to be in slice: %s.",
		productName,
		mainVersions,
	)}
}

func (p *publisher) publishVersionAssets(pr productRoot, version string) error {
	vPath := pr.versionFilePath(version)
	assetsDir := GetAssetsDir()
	imgDirName := "images"

	if assetsDir == "" || assetsDir == getDocsDir() {
		log(lInfo, "Assets directory is not configured or is same as docs dir. Skipping asset publication for `%s` version `%s`.\n", pr.Key, version)
		return nil
	}

	targetDir := fmt.Sprintf("%s%c%s%c%s", assetsDir, os.PathSeparator, pr.Key, os.PathSeparator, version)
	log(lInfo, "Publishing assets for version `%s`. Target dir: `%s`\n", version, targetDir)

	// publish images
	imgSrcDir := fmt.Sprintf("%s%c%s", vPath, os.PathSeparator, imgDirName)
	imgTargetDir := fmt.Sprintf("%s%c%s", targetDir, os.PathSeparator, imgDirName)

	if _, err := os.Stat(imgSrcDir); !os.IsNotExist(err) {
		if err := cp.Copy(imgSrcDir, imgTargetDir); err != nil {
			return err
		}
	}

	return nil
}
