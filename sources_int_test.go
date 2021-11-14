//go:build integration || ci

package docweaver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadSources(t *testing.T) {
	_ = os.Setenv(EnvKeySourcesFile, "./testdata/doc-sources.yml")
	expectedSources := []source{
		{
			Key: "scavenger",
			Url: "https://github.com/reliqarts/scavenger-docs",
		},
		{
			Key: "docweaver",
			Url: "https://github.com/reliqarts/docweaver-docs",
		},
	}
	sources, err := readSources()
	if err != nil {
		t.Fatal(err)
	}

	for _, es := range expectedSources {
		assert.Contains(t, fmt.Sprintf("%s", sources), fmt.Sprintf("%s %s", es.Key, es.Url))
	}

	assert.NotContains(t, fmt.Sprintf("%s", sources), "sources")
}
