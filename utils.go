package docweaver

import (
	"fmt"
	"github.com/reliqarts/go-common"
	"golang.org/x/net/html"
	"log"
	"os"
	"sort"
	"strconv"
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
	return common.GetEnvOrDefault(EnvKeyAssetsDir, getDocsDir())
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

func getDocsDir() string {
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

// sortVersions sorts a given slice of versions
func sortVersions(versions []string) {
	const highVal = 99999.99
	sort.SliceStable(versions, func(i, j int) bool {
		verI := strings.Split(versions[i], "-")[0]
		verJ := strings.Split(versions[j], "-")[0]

		if verI[0] == 'v' || verI[0] == 'V' {
			verI = verI[1:]
		}
		if verJ[0] == 'v' || verJ[0] == 'V' {
			verJ = verJ[1:]
		}

		jv, err := strconv.ParseFloat(verI, 32)
		if err != nil {
			jv = highVal
		}
		iv, err := strconv.ParseFloat(verJ, 32)
		if err != nil {
			iv = highVal
		}
		return jv < iv
	})
}

func latestVersion(versions []string) (latest string) {
	var vs []string

	// focus on non-main versions
	for _, v := range versions {
		if v != versionMaster && v != versionMain {
			vs = append(vs, v)
		}
	}
	if len(vs) == 0 {
		// no non-main versions exist, i.e. no version
		return "N/A"
	}

	sortVersions(vs)
	for i := len(vs) - 1; i >= 0; i-- {
		verVal := vs[i]
		if verVal[0] == 'v' || verVal[0] == 'V' {
			verVal = verVal[1:]
		}
		_, err := strconv.ParseFloat(verVal, 32)
		if err == nil {
			latest = vs[i]
			break
		}
	}

	if latest == "" {
		// pick the latest by alpha sort
		sort.Strings(vs)
		latest = vs[len(vs)-1]
	}

	return
}
