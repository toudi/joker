package utils_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/toudi/joker/internal/utils"
)

func TestPathToPathWithPattern(t *testing.T) {
	type testCase struct {
		input    string
		expected string
	}

	for _, test := range []testCase{
		{
			input:    "foo/bar/file.txt",
			expected: "foo/bar/file.txt",
		},
		{
			input:    "/foo/bar",
			expected: "/foo/bar/*",
		},
		{
			input:    "/foo/bar/*",
			expected: "/foo/bar/*",
		},
		{
			input:    "/foo/bar/*.go",
			expected: "/foo/bar/*.go",
		},
	} {
		assert.Equal(t, test.expected, utils.PathToPathWithPattern(test.input))
	}
}
