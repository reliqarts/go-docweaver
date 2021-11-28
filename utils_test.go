//go:build unit || ci

package docweaver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStringSet struct {
	key      string
	input    []string
	expected []string
}

func TestGetLoggerSet(t *testing.T) {
	loggerSet := GetLoggerSet()

	assert.Contains(t, loggerSet.Info.Prefix(), "info")
	assert.Contains(t, loggerSet.Err.Prefix(), "err")
	assert.Contains(t, loggerSet.Warn.Prefix(), "warn")
}

func TestReplaceLinks(t *testing.T) {
	productKey := "prod-up"
	testData := []struct {
		version  string
		content  string
		expected string
	}{
		{
			version:  "2.0",
			content:  "{{docs}}/something/somewhere",
			expected: fmt.Sprintf("%s/%s/2.0/something/somewhere", GetRoutePrefix(), productKey),
		},
		{
			version:  "2.0",
			content:  "docs/{{version}}/something/somewhere",
			expected: fmt.Sprintf("%s/%s/2.0/something/somewhere", GetRoutePrefix(), productKey),
		},
		{
			version:  "v5.0",
			content:  "{{version}}/some-thing/somewhere",
			expected: "v5.0/some-thing/somewhere",
		},
		{
			version:  "6.0-alpha",
			content:  "{{version}}/some-thing/somewhere",
			expected: "6.0-alpha/some-thing/somewhere",
		},
	}

	for _, td := range testData {
		assert.Equal(t, td.expected, replaceLinks(productKey, td.version, td.content))
	}
}

func TestSortVersions(t *testing.T) {
	testData := []testStringSet{
		{
			key:      "1",
			input:    []string{"1.0", "2.0-beta", "master", "main", "10.0", "20", "2"},
			expected: []string{"1.0", "2.0-beta", "2", "10.0", "20", "master", "main"},
		},
		{
			key:      "2",
			input:    []string{"1.0", "3.0-beta", "main", "10.0", "20.5", "200"},
			expected: []string{"1.0", "3.0-beta", "10.0", "20.5", "200", "main"},
		},
		{
			key:      "3",
			input:    []string{"1.0", "3.0-beta", "main", "10.0", "0.0200"},
			expected: []string{"0.0200", "1.0", "3.0-beta", "10.0", "main"},
		},
	}

	for _, td := range testData {
		t.Run(td.key, func(t *testing.T) {
			sortVersions(td.input)

			assert.Equal(t, fmt.Sprintf("%s", td.expected), fmt.Sprintf("%s", td.input))
		})
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
		t.Run(td.expected, func(t *testing.T) {
			result := latestVersion(td.input)

			assert.Equal(t, td.expected, result)
		})
	}
}
