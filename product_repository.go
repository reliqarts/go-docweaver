package docweaver

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Cleaner interface {
	CleanTempVersions() error
}

type ProductRepository interface {
	Cleaner
	FindAllProducts() ([]Product, error)
	FindProduct(productName string) (*Product, error)
	GetDir() string
	GetPage(productName, version, pagePath string) (*Page, error)
	GetIndex(productName string) (*Page, error)
	ListProductKeys() ([]string, error)
}

type productRepository struct {
	dir string
}

const (
	defaultPagePath = "installation"
	pageExt         = "md"
	indexPath       = "documentation"
)

func GetRepository(dir string) ProductRepository {
	if dir == "" {
		dir = getDocsDir()
	}
	return &productRepository{dir: dir}
}

func (pr *productRepository) GetDir() string {
	return pr.dir
}

// FindAllProducts finds all products from the documentations' directory.
func (pr *productRepository) FindAllProducts() ([]Product, error) {
	pn, err := pr.ListProductKeys()
	if err != nil {
		return nil, err
	}

	var products []Product
	for _, n := range pn {
		p, err := pr.FindProduct(n)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

// ListProductKeys lists keys of all available products.
func (pr *productRepository) ListProductKeys() ([]string, error) {
	var productNames []string
	cmd := exec.Command("ls")
	cmd.Dir = pr.dir

	out, err := cmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log(lWarn, "Docs directory `%s` does not exist.\n", pr.dir)
			return productNames, nil
		}

		return productNames, simpleError{fmt.Sprintf("Failed to list all products in docs dir: ``. %s\n", err)}
	}

	for _, pn := range strings.Split(string(out), "\n") {
		n := strings.TrimSpace(pn)
		if n != "" {
			productNames = append(productNames, n)
		}
	}

	return productNames, nil
}

func (pr *productRepository) FindProduct(productKey string) (*Product, error) {
	r := productRoot{ParentDir: pr.dir, Key: productKey}
	var versions []string

	entries, err := os.ReadDir(r.filePath())
	if err != nil {
		return nil, err
	}

	for _, f := range entries {
		if f.IsDir() {
			vt := f.Name()
			if strings.Contains(vt, tempNameSuffix) {
				continue
			}
			versions = append(versions, vt)
		}
	}

	return pr.newProduct(r, versions), nil
}

func (pr *productRepository) GetPage(productKey, version, pagePath string) (*Page, error) {
	if productKey == "" {
		err := simpleError{err: "No product key provided."}
		log(lError, err.Error())
		return nil, err
	}
	if version == "" {
		log(lInfo, "Using default version (%s) for product `%s`, page path: `%s`.\n", defaultVersion, productKey, pagePath)
		version = defaultVersion
	}
	if pagePath == "" {
		log(lInfo, "Using default page path (%s) for product `%s`, version: `%s`.\n", defaultPagePath, productKey, version)
		pagePath = defaultPagePath
	}

	r := productRoot{ParentDir: pr.dir, Key: productKey}
	p, err := pr.FindProduct(productKey)
	if err != nil {
		return nil, simpleError{fmt.Sprintf("Failed to init product for page. %s", err)}
	}

	filePath := fmt.Sprintf("%s%c%s.%s", r.versionFilePath(version), os.PathSeparator, pagePath, pageExt)
	md, err := os.ReadFile(filePath)
	if err != nil {
		log(lWarn, "Failed to read product page from file path `%s`.\n", filePath)
		return nil, err
	}

	var rawContent bytes.Buffer
	err = goldmark.New(
		goldmark.WithExtensions(extension.GFM, emoji.Emoji),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps(), html.WithXHTML(), html.WithUnsafe()),
	).Convert(md, &rawContent)
	if err != nil {
		return nil, err
	}

	content := replaceLinks(productKey, version, rawContent.String())

	var index *Page = nil
	if pagePath != indexPath {
		index, err = pr.GetPage(r.Key, version, indexPath)
		if err != nil {
			log(lWarn, "Failed to read product index for page (%s) from path `%s`.\n", pagePath, indexPath)
			index = nil
		}
	}

	return &Page{
		UrlPath: pagePath,
		Title:   getPageTitleFromHtml(content),
		Content: content,
		Product: p,
		Version: version,
		Index:   index,
	}, nil
}

func (pr *productRepository) GetIndex(productName string) (*Page, error) {
	return pr.GetPage(productName, defaultVersion, defaultPagePath)
}

// CleanTempVersions removes all temporary documentation versions. Only returns the last error that occurred.
func (pr *productRepository) CleanTempVersions() (lastErr error) {
	products, err := pr.FindAllProducts()
	if err != nil {
		lastErr = err
		return
	}
	for _, p := range products {
		for _, ver := range p.Versions {
			lastErr = removeDir(p.root.versionFilePath(versionTempName(ver)))
		}
	}
	return
}

func (pr *productRepository) newProduct(r productRoot, versions []string) (product *Product) {
	latestV := latestVersion(versions)
	product = &Product{
		Name:          cases.Title(language.English).String(r.Key),
		BaseUrl:       fmt.Sprintf("%s/%s", GetRoutePrefix(), r.Key),
		LatestVersion: latestV,
		Versions:      versions,
		root:          r,
	}
	product.loadMeta()
	return
}
