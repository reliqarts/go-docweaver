package docweaver

import (
	"fmt"
	"github.com/reliqarts/go-common"
	"os"
	"os/exec"
	"strings"
)

const (
	defaultDocumentationDir string = "./tmp/docs"
	versionMaster           tag    = "master"
	versionMain             tag    = "main"
)

type Publisher interface {
	Publish(productName string, source string, shouldUpdate bool)
	GetDocsDir() string
}

type publisher struct {
	docsDir string
}

// GetPublisher returns the default instance of Publisher.
func GetPublisher() Publisher {
	return publisher{
		docsDir: common.GetEnvOrDefault(envKeyName("docs_dir"), defaultDocumentationDir),
	}
}

// GetPublisherWithDocsDir returns an instance of Publisher with the provided [docsDir].
func GetPublisherWithDocsDir(docsDir string) Publisher {
	return publisher{docsDir: docsDir}
}

func (p publisher) Publish(productName string, source string, shouldUpdate bool) {
	pRoot := productRoot{ParentDir: p.docsDir, Name: productName, Source: source}
	mainVersions := []tag{versionMaster, versionMain}
	baseVersion := tag("")

	loggers.Info.Println("Product root:", pRoot)

	for _, mv := range mainVersions {
		if err := p.publishProductVersion(pRoot, mv, true); err == nil {
			baseVersion = mv
			break
		}
	}

	if baseVersion == "" {
		loggers.Err.Fatalf(
			"Base version for product %s could not be determined. Was not found to be in slice: %s.",
			productName,
			mainVersions,
		)
	}

	tags, err := p.listProductTags(pRoot, baseVersion)
	if err != nil {
		loggers.Err.Printf("Failed to list tags. %s\n", err)
	}

	for _, tag := range tags {
		if err := p.publishProductVersion(pRoot, tag, shouldUpdate); err != nil {
			loggers.Warn.Printf("Failed to publish/shouldUpdate tag `%s`.", tag)
		}
	}
}

func (p publisher) GetDocsDir() string {
	return p.docsDir
}

func (p *publisher) publishProductVersion(pr productRoot, version tag, update bool) error {
	prFullPath := pr.fullPath()
	verPath := fmt.Sprintf("%s/%s", prFullPath, version)
	tempNameSuffix := common.RandomString(6)
	verPathTemp := fmt.Sprintf(verPath + "-" + tempNameSuffix)
	tempVerPathExists := false

	if _, err := os.Stat(prFullPath); os.IsNotExist(err) {
		loggers.Info.Printf("Project root not present, creating project root: `%s`\n", prFullPath)
		if err = os.MkdirAll(prFullPath, 0755); err != nil {
			return err
		}
	}

	if _, err := os.Stat(verPath); !os.IsNotExist(err) {
		if !update {
			loggers.Info.Printf("Version `%s` already exists for product `%s`. Update not requested. Skipped.\n", version, pr.Name)
			return nil
		}

		loggers.Info.Printf("Path for version `%s` already existed. Removing `%s` for update.\n", version, verPath)
		// temporarily rename existing version
		if err := os.Rename(verPath, verPathTemp); err != nil {
			loggers.Err.Printf("Failed to temporarily rename version path `%s` to `%s` for update. %s\n", verPath, verPathTemp, err)
			return err
		}
		tempVerPathExists = true
	}

	loggers.Info.Printf("Executing version clone `%s` into: `%s`\n", version, prFullPath)
	cmd := exec.Command("git", "clone", "--branch", string(version), pr.Source, string(version))
	cmd.Dir = pr.fullPath()
	if _, err := cmd.Output(); err != nil {
		loggers.Err.Printf("Failed to execute version clone `%s` into: `%s`. %s\n", version, prFullPath, err)
		return err
	}

	// remove temp. version path if exists
	if tempVerPathExists {
		if err := os.RemoveAll(verPathTemp); err != nil {
			loggers.Err.Printf("Failed to remove temporary version path `%s` after update. %s\n", verPathTemp, err)
			return err
		}
	}

	return nil
}

func (p *publisher) listProductTags(pr productRoot, baseVersion tag) ([]tag, error) {
	var tags []tag
	loggers.Info.Printf("Listing tags for product `%s` using base version `%s`.\n", pr.Name, baseVersion)
	cmd := exec.Command("git", "tag", "-l")
	cmd.Dir = fmt.Sprintf("%s/%s", pr.fullPath(), baseVersion)
	out, err := cmd.Output()
	if err != nil {
		loggers.Err.Printf("Failed to list tags for product `%s` using base version `%s`.\n", pr.Name, baseVersion)
		return nil, err
	}

	for _, t := range strings.Split(string(out), "\n") {
		if t != "" {
			tags = append(tags, tag(t))
		}
	}

	return tags, nil
}
