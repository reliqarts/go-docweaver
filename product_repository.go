package docweaver

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type ProductRepository interface {
	GetDir() string
	FindAllProducts() ([]Product, error)
	FindProduct(productName string) (*Product, error)
	ListProductKeys() ([]string, error)
	GetPage(productName, version, pagePath string) (*Page, error)
	GetIndex(productName string) (*Page, error)
}

type productRepository struct {
	dir string
}

const (
	defaultPagePath = "installation"
	pageExt         = "md"
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

	files, err := ioutil.ReadDir(r.filePath())
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			vt := f.Name()
			versions = append(versions, vt)
		}
	}

	return pr.newProduct(r, versions), nil
}

func (pr *productRepository) newProduct(r productRoot, versions []string) (product *Product) {
	product = &Product{
		Name:          strings.Title(r.Key),
		Versions:      versions,
		LatestVersion: latestVersion(versions),
		root:          r,
	}
	product.loadMeta()
	return
}

func (pr *productRepository) GetPage(productKey, version, pagePath string) (*Page, error) {
	r := productRoot{ParentDir: pr.dir, Key: productKey}

	if productKey == "" {
		err := simpleError{err: "No product key provided."}
		loggers.Err.Println(err)
		return nil, err
	}
	if version == "" {
		loggers.Info.Printf("Using default version (%s) for product `%s`, page path: `%s`.\n", defaultVersion, productKey, pagePath)
		version = defaultVersion
	}
	if pagePath == "" {
		loggers.Info.Printf("Using default page path (%s) for product `%s`, version: `%s`.\n", defaultPagePath, productKey, version)
		pagePath = defaultPagePath
	}

	filePath := fmt.Sprintf("%s%c%s.%s", r.versionFilePath(version), os.PathSeparator, pagePath, pageExt)
	md, err := os.ReadFile(filePath)
	if err != nil {
		loggers.Err.Printf("Failed to read product page from file path `%s`.\n", filePath)
		return nil, err
	}

	content := replaceLinks(productKey, version, string(markdown.ToHTML(md, nil, nil)))

	return &Page{
		UrlPath: pagePath,
		Title:   getPageTitleFromHtml(content),
		Content: content,
		Version: version,
	}, nil
}

func (pr *productRepository) GetIndex(productName string) (*Page, error) {
	return pr.GetPage(productName, defaultVersion, defaultPagePath)
}
