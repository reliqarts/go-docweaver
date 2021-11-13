//go:build unit

package docweaver_test

import (
	"github.com/reliqarts/go-docweaver"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	defaultDocsDir string = "./tmp/docs"
)

func TestGetPublisher(t *testing.T) {
	pub := docweaver.GetPublisher()

	assert.Equal(t, defaultDocsDir, pub.GetDocsDir())
}

func TestGetPublisherWithDocsDir(t *testing.T) {
	dir := "foo"
	pub := docweaver.GetPublisherWithDocsDir(dir)

	assert.Equal(t, dir, pub.GetDocsDir())
}