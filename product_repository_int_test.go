//go:build integration || ci

package docweaver_test

import (
	"fmt"
	"github.com/reliqarts/go-docweaver"
	"github.com/stretchr/testify/assert"
	"testing"
)

var repo = docweaver.GetRepository(docsDir)

func TestProductRepository_FindProduct(t *testing.T) {
	expectedName := "Product One"
	product, err := repo.FindProduct(testProductKey)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Product:", product)

	assert.Equal(t, expectedName, product.Name)
	assert.Equal(t, "Simple test product.", product.Description)
	assert.Contains(t, product.ImageUrl, fmt.Sprintf("%s/%s/main/", docweaver.GetAssetsRoutePrefix(), testProductKey))
}

func TestProductRepository_GetPage(t *testing.T) {
	versions := []string{"main", "1.0"}

	for _, version := range versions {
		page, err := repo.GetPage(testProductKey, version, "installation")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("%v\n", page)

		assert.Equal(t, "Product 1", page.Title)
		assert.Equal(t, version, page.Version)
		assert.Contains(t, page.Content, fmt.Sprintf("docs/%s/%s/", testProductKey, version))
	}
}
