//go:build integration || ci

package docweaver_test

import (
	"fmt"
	"github.com/reliqarts/go-docweaver"
	"os"
	"testing"
)

const (
	docsDir        string = "./testdata/docs"
	testProductKey string = "product1"
)

func init() {
	_ = os.Setenv(docweaver.EnvKeyAssetsDir, "./testdata/assets")
}

func TestPublisher_Publish(t *testing.T) {
	productName := "docweaver"
	productPath := docsDir + "/" + productName
	publisher := docweaver.GetPublisherWithDocsDir(docsDir)
	publisher.Publish(productName, "https://github.com/reliqarts/docweaver-docs.git", true)
	versionsToCheck := []string{"main", "1.0", "2.0", "3.0", "4.0"}

	if _, err := os.Stat(productPath); os.IsNotExist(err) {
		t.Fatalf("Product directory not found for product `%s`", productName)
	}

	for _, v := range versionsToCheck {
		if _, err := os.Stat(productPath + "/" + v); os.IsNotExist(err) {
			t.Fatalf("Product version not found `%s`", productPath+"/"+v)
		}
	}

	if err := os.RemoveAll(productPath); err != nil {
		t.Fatalf("Failed to remove product path `%s`. %s", productPath, err)
	}

	assetsDir := docweaver.GetAssetsDir()
	if assetsDir != "" {
		assetsPath := fmt.Sprintf("%s%c%s", assetsDir, os.PathSeparator, productName)
		if err := os.RemoveAll(assetsPath); err != nil {
			t.Fatalf("Failed to remove published product assets at path: `%s`. %s", assetsPath, err)
		}
	}
}
