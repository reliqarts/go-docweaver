# Docweaver

An easy-to-use product documentation package in Golang.

Docweaver is suitable for product documentation and/or knowledge bases. Converts folder(s) of .md files into full-bread
complete documentation. This package is without UI elements and templates, all freedom regarding final presentation is
given to the end-user.

PHP/Laravel version is available [here](https://github.com/reliqarts/laravel-docweaver).

[![Go Reference](https://pkg.go.dev/badge/github.com/reliqarts/go-docweaver.svg)](https://pkg.go.dev/github.com/reliqarts/go-docweaver)
[![Build Status](https://github.com/reliqarts/go-docweaver/workflows/CI/badge.svg)](https://github.com/reliqarts/go-docweaver/actions?query=workflow:CI)
[![Codecov](https://img.shields.io/codecov/c/github/reliqarts/go-docweaver.svg?style=flat)](https://codecov.io/gh/reliqarts/go-docweaver)
[![https://goreportcard.com/report/github.com/reliqarts/go-docweaver](https://goreportcard.com/badge/github.com/reliqarts/go-docweaver)](https://goreportcard.com/report/github.com/reliqarts/go-docweaver)

---

## Installation & Usage

### Installation

```bash
go get github.com/reliqarts/go-docweaver
```

### Setup

Example .env:

```dotenv
DW_DOCS_DIR=./tmp/docs               # Where documentation repos (archives) should be stored.
DW_ASSETS_DIR=./tmp/doc-assets       # Where documentation assets should be accessed from.
DW_ROUTE_PREFIX=docs                 # Documentation route prefix.
DW_ASSETS_ROUTE_PREFIX=doc-assets    # Route prefix for assets.
DW_SOURCES_FILE=./doc-sources.yml    # Sources file location.
```

Example files:
- [doc-sources.yml](https://github.com/reliqarts/go-docweaver/blob/main/testdata/doc-sources.yml)

#### Documentation Directory

The documentation directory is the place where you put your project documentation directories. It may be changed with
the environment variable `DW_DOCS_DIR`. The default documentation directory is `./tmp/docs`.

#### Structure

Each project directory should contain separate folders for each documented version. Each version must have at least
two (2) markdown files, namely `documentation.md` and `installation.md`, which serve as the index (usually shown in sidebar) and initial
documentation pages respectively.

```
[doc dir]
    │
    └─── Project One
    │       └── 1.0 
    │       └── 2.1
    │            └── .docweaver.yml       # meta file (optional)
    │            └── documentation.md     # sidebar nav
    │            └── installation.md      # initial page
    │
    └─── Project Two
```

#### Meta File

Configurations for each doc version may be placed in `.docweaver.yml`. The supported settings are:

- #### name
  Product name.
- #### description
  Product description.
- #### image_url
  Product image url. This may be an absolute url (e.g. `http://mywebsite.com/myimage.jpg`) or an image found in
  the `images` resource directory.

  To use the `foo.jpg` in the `images` directory you would set `image_url` to `{{docs}}/images/foo.jpg`.


### Usage

<details>
<summary>Gin Example</summary>

###### main.go

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/reliqarts/go-docweaver"
)

func main() {
	router := gin.New()

	router.GET(docweaver.GetRoutePrefix(), handlers.Documentation())
	router.GET(fmt.Sprintf("%s/*path", docweaver.GetRoutePrefix()), handlers.Documentation())
	router.Static(docweaver.GetAssetsRoutePrefix(), docweaver.GetAssetsDir())

	_ = (router).Run("localhost:8080")
}

// ...
```

###### handlers.go

```go
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/reliqarts/go-docweaver"
	"net/http"
)

// ...

func Documentation() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Param("path")
		dw := docweaver.GetRepository("")

		if path == "/" || path == "" {
			products, err := dw.FindAllProducts()
			if err != nil {
				c.HTML(code, "error.gohtml", gin.H{
					"errorCode":    http.StatusInternalServerError,
					"errorMessage": err,
				})
				return
			}

			c.HTML(http.StatusOK, "docs/index.gohtml", gin.H{
				"products": products,
			})
			return
		}

		productKey, version, pagePath := "", "", ""
		pageParts := strings.Split(path, "/")
		if len(pageParts) >= 2 {
			productKey = pageParts[1]
		}
		if len(pageParts) >= 3 {
			version = pageParts[2]
		}
		if len(pageParts) >= 4 {
			pagePath = pageParts[3]
		}

		page, err := dw.GetPage(productKey, version, pagePath)
		if err != nil {
			errMsg := fmt.Sprintf("Page not found. %s", err)
			c.HTML(http.StatusNotFound, "error.gohtml", gin.H{
				"errorCode":    http.StatusNotFound,
				"errorMessage": errMsg,
			})
			c.Abort()
			return
		}

		c.HTML(http.StatusOK, "documentation/show.gohtml", gin.H{
			"page": page,
		})
	}
}

// ...
```

</details>

&nbsp;

:beers: cheers
