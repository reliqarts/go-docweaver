//go:build unit || ci

package docweaver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStringSet struct {
	input    []string
	expected []string
}

func TestGetLoggerSet(t *testing.T) {
	loggerSet := GetLoggerSet()

	assert.Contains(t, loggerSet.Info.Prefix(), "info")
	assert.Contains(t, loggerSet.Err.Prefix(), "err")
	assert.Contains(t, loggerSet.Warn.Prefix(), "warn")
}

func TestSortVersions(t *testing.T) {
	testData := []testStringSet{
		{
			input:    []string{"1.0", "2.0-beta", "master", "main", "10.0", "20", "2"},
			expected: []string{"1.0", "2.0-beta", "2", "10.0", "20", "master", "main"},
		},
		{
			input:    []string{"1.0", "3.0-beta", "main", "10.0", "20.5", "200"},
			expected: []string{"1.0", "3.0-beta", "10.0", "20.5", "200", "main"},
		},
		{
			input:    []string{"1.0", "3.0-beta", "main", "10.0", "0.0200"},
			expected: []string{"0.0200", "1.0", "3.0-beta", "10.0", "main"},
		},
	}

	for _, td := range testData {
		sortVersions(td.input)

		assert.Equal(t, fmt.Sprintf("%s", td.expected), fmt.Sprintf("%s", td.input))
	}
}

func TestLatestVersion(t *testing.T) {
	testData := []struct {
		input    []string
		expected string
	}{
		{
			[]string{"1.0", "3.0-beta", "main", "10.0", "0.0200"},
			"10.0",
		},
		{
			[]string{"feature/foo", "v1.5", "1.0", "1.04", "master", "develop"},
			"v1.5",
		},
		{
			[]string{"0.1-beta", "0.1-alpha", "main"},
			"0.1-beta",
		},
		{
			[]string{"master", "main"},
			"N/A",
		},
		{
			[]string{},
			"N/A",
		},
	}

	for _, td := range testData {
		result := latestVersion(td.input)

		assert.Equal(t, td.expected, result)
	}
}
