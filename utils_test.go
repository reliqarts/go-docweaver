//go:build unit || ci

package docweaver_test

import (
	"github.com/reliqarts/go-docweaver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLoggerSet(t *testing.T) {
	loggerSet := docweaver.GetLoggerSet()

	assert.Contains(t, loggerSet.Info.Prefix(), "info")
	assert.Contains(t, loggerSet.Err.Prefix(), "err")
	assert.Contains(t, loggerSet.Warn.Prefix(), "warn")
}
