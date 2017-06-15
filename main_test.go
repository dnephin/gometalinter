package main

import (
	"io/ioutil"
	"sort"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

func TestRelativePackagePath(t *testing.T) {
	var testcases = []struct{
		dir string
		expected string
	}{
		{
			dir: "/abs/path",
			expected: "/abs/path",
		},
		{
			dir: ".",
			expected: ".",
		},
		{
			dir: "./foo",
			expected: "./foo",
		},
		{
			dir: "relative/path",
			expected: "./relative/path",
		},
	}

	for _, testcase := range testcases {
		assert.Equal(t, testcase.expected, relativePackagePath(testcase.dir))
	}
}

func TestSortedIssues(t *testing.T) {
	actual := []*Issue{
		{Path: "b.go", Line: 5},
		{Path: "a.go", Line: 3},
		{Path: "b.go", Line: 1},
		{Path: "a.go", Line: 1},
	}
	issues := &sortedIssues{
		issues: actual,
		order:  []string{"path", "line"},
	}
	sort.Sort(issues)
	expected := []*Issue{
		{Path: "a.go", Line: 1},
		{Path: "a.go", Line: 3},
		{Path: "b.go", Line: 1},
		{Path: "b.go", Line: 5},
	}
	require.Equal(t, expected, actual)
}

func TestLoadConfigWithDeadline(t *testing.T) {
	originalConfig := *config
	defer func() { config = &originalConfig }()

	tmpfile, err := ioutil.TempFile("", "test-config")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(`{"Deadline": "3m"}`))
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	filename := tmpfile.Name()
	err = loadConfig(nil, &kingpin.ParseElement{Value: &filename}, nil)
	require.NoError(t, err)

	require.Equal(t, 3 * time.Minute, config.Deadline.Duration())
}