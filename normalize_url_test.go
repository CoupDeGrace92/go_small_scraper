package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
		err      bool
	}{
		{
			name:     "remove scheme",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
			err:      false,
		},
		{
			name:     "remove trailing backslash",
			inputURL: "https://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
			err:      false,
		},
		{
			name:     "upper to lower host",
			inputURL: "https://WWW.BOOT.DEV/blog/path",
			expected: "www.boot.dev/blog/path",
			err:      false,
		},
		{
			name:     "malformed url",
			inputURL: "This is not a url",
			expected: "",
			err:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
