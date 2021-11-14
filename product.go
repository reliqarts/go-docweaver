package docweaver

import (
	"fmt"
	yml "gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Product struct {
	Name        string
	Description string
	ImageUrl    string
	Versions    []string
	root        productRoot
}

type Page struct {
	UrlPath string
	Title   string
	Content string
	Version string
}

type productRoot struct {
	ParentDir string
	Key       string
	Source    string
}

type productMeta struct {
	Name        string
	Description string
	ImageUrl    string `yaml:"image_url"`
}

func (p *productRoot) filePath() string {
	return fmt.Sprintf("%s%c%s", p.ParentDir, os.PathSeparator, p.Key)
}

func (p *productRoot) versionFilePath(version string) string {
	return fmt.Sprintf("%s%c%s", p.filePath(), os.PathSeparator, version)
}

func (p *productRoot) hasSource() bool {
	return p.Source != ""
}

func (p *Product) GetKey() string {
	return p.root.Key
}

func (p *Product) loadMeta() {
	var err error
	var meta *productMeta
	var mainVersion string
	r := p.root

	// find product meta using main versions
	for _, ver := range intersection(p.Versions, mainVersions) {
		meta, err = p.readMeta(ver)
		if err != nil {
			loggers.Err.Printf("Failed to read meta file from product `%s`, version `%s`. %s\n", r.Key, ver, err)
		}
		if meta != nil {
			mainVersion = ver
			break
		}
	}

	if meta != nil {
		if meta.Name != "" {
			p.Name = meta.Name
		}
		if meta.Description != "" {
			p.Description = meta.Description
		}
		if meta.ImageUrl != "" {
			p.ImageUrl = p.getAssetLink(meta.ImageUrl, mainVersion)
		}
	}
}

func (p *Product) readMeta(version string) (meta *productMeta, err error) {
	r := p.root
	metaFilePath := fmt.Sprintf("%s%c%s", r.versionFilePath(version), os.PathSeparator, metaFileName)

	yaml, err := os.ReadFile(metaFilePath)
	if err != nil {
		return nil, err
	}

	if err := yml.Unmarshal(yaml, &meta); err != nil {
		return nil, err
	}

	return meta, nil
}

func (p *Product) getAssetLink(link, version string) string {
	linkReplacement := fmt.Sprintf("%s/%s/%s", GetAssetsRoutePrefix(), p.root.Key, version)
	repl := strings.NewReplacer(assetUrlPlaceholder, linkReplacement)

	return repl.Replace(link)
}
