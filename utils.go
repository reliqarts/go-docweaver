package docweaver

import (
	"fmt"
	"github.com/reliqarts/go-common"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"
)

type loggerSet struct {
	Err  *log.Logger
	Info *log.Logger
	Warn *log.Logger
}

type simpleError struct {
	err string
}

const (
	EnvKeyDocsDir           string = "DW_DOCS_DIR"            // Docs directory environment key.
	EnvKeyAssetsDir         string = "DW_ASSETS_DIR"          // Assets directory environment key.
	EnvKeyRoutePrefix       string = "DW_ROUTE_PREFIX"        // Route prefix environment key.
	EnvKeyAssetsRoutePrefix string = "DW_ASSETS_ROUTE_PREFIX" // Assets route prefix environment key.
	EnvKeySourcesFile       string = "DW_SOURCES_FILE"        // Sources file environment key.

	defaultDocumentationDir  string = "./tmp/docs"
	defaultVersion                  = versionMain
	defaultRoutePrefix              = "docs"
	defaultAssetsRoutePrefix        = "doc-assets"
	defaultSourcesFile              = "./doc-sources.yml"

	metaFileName = ".docweaver.yml"

	versionMaster       string = "master"
	versionMain         string = "main"
	versionPlaceholder  string = "{{version}}"
	assetUrlPlaceholder string = "{{docs}}"
)

var loggers = GetLoggerSet()

func (e simpleError) Error() string {
	return e.err
}

// GetLoggerSet returns configured loggers for package.
func GetLoggerSet() *loggerSet {
	return &loggerSet{
		Err:  log.New(os.Stdout, "[Dw][err] ", log.Ldate|log.Ltime),
		Info: log.New(os.Stdout, "[Dw][info] ", log.Ldate|log.Ltime),
		Warn: log.New(os.Stdout, "[Dw][warn] ", log.Ldate|log.Ltime),
	}
}

// GetAssetsDir returns configured assets directory. env key: DW_ASSETS_DIR
func GetAssetsDir() string {
	return common.GetEnvOrDefault(EnvKeyAssetsDir, "")
}

// GetRoutePrefix returns configured documentation route prefix. env key: DW_ROUTE_PREFIX
func GetRoutePrefix() string {
	return common.GetEnvOrDefault(EnvKeyRoutePrefix, defaultRoutePrefix)
}

// GetAssetsRoutePrefix returns configured assets route prefix. env key: DW_ASSETS_ROUTE_PREFIX
func GetAssetsRoutePrefix() string {
	return common.GetEnvOrDefault(EnvKeyAssetsRoutePrefix, defaultAssetsRoutePrefix)
}

// GetSourcesFilePath returns configured sources file path. env key: DW_SOURCES_FILE
func GetSourcesFilePath() string {
	return common.GetEnvOrDefault(EnvKeySourcesFile, defaultSourcesFile)
}

func getConfiguredDocsDir() string {
	return common.GetEnvOrDefault(EnvKeyDocsDir, defaultDocumentationDir)
}

func getPageTitleFromHtml(content string) string {
	z := html.NewTokenizer(strings.NewReader(content))
	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			tag := z.Token()

			if tag.Data == "h1" {
				if tt = z.Next(); tt == html.TextToken {
					return z.Token().Data
				}
			}
		}
	}

	return ""
}

func replaceLinks(productKey, version, content string) string {
	linkReplacement := fmt.Sprintf("%s/%s/%s", GetRoutePrefix(), productKey, version)
	repl := strings.NewReplacer(
		assetUrlPlaceholder, linkReplacement,
		fmt.Sprintf("%s%s", "docs/", versionPlaceholder), linkReplacement,
		versionPlaceholder, version,
	)

	return repl.Replace(content)
}

func intersection(s1, s2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove duplicates from slice.
	inter = removeDuplicates(inter)
	return
}

//Remove duplicates from slice.
func removeDuplicates(elements []string) (nodups []string) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}
