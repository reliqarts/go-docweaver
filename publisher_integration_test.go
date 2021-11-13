//go:build integration

package docweaver_test

import (
	"github.com/reliqarts/go-docweaver"
	"os"
	"testing"
)

const testDocsDir string = "./testdata/docs"

func TestPublisher_Publish(t *testing.T) {
	docsDir := testDocsDir
	productName := "docweaver"
	productPath := docsDir + "/" + productName
	publisher := docweaver.GetPublisherWithDocsDir(docsDir)
	publisher.Publish(productName, "https://github.com/reliqarts/docweaver-docs.git", true)
	versionsToCheck := []string{"master", "1.0", "2.0", "3.0", "4.0"}

	if _, err := os.Stat(productPath); os.IsNotExist(err) {
		t.Fatalf("Product directory not found for product `%s`", productName)
	}

	for _, v := range versionsToCheck {
		if _, err := os.Stat(productPath + "/" + v); os.IsNotExist(err) {
			t.Fatalf("Product version not found `%s`", productPath+"/"+v)
		}
	}

	if err := os.RemoveAll(productPath); err != nil {
		t.Fatalf("Failed to cleanup after test. %s", err)
	}
}
