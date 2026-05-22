package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
		err      bool //in this case, until we have test cases that can err, this should be false
	}{
		{
			name: "Base case",
			body: `
			<html>
				<body>
					<h1>Test title</h1>
					<h2>Sub-title</h2>
				</body>
			</html>
			`,
			expected: "Test title",
			err:      false,
		},
		{
			name: "No h1",
			body: `
			<html>
				<body>
					<h2>Sub-title</h2>
				</body>
			</html>
			`,
			expected: "Sub-title",
			err:      false,
		},
		{
			name: "No headers",
			body: `
			<html>
				<body>

				</body>
			</html>
			`,
			expected: "",
			err:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			heading, err := getHeadingFromHTML(tc.body)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, heading)
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
		err      bool //currently should only ever be false until we have tcs that can err
	}{
		{
			name: "Base case",
			body: `
			<html><body>
				<p>Outside paragraph</p>
				<main>
					<p>Inside paragraph</p>
				</main>
			</body></html>
			`,
			expected: "Inside paragraph",
			err:      false,
		},
		{
			name: "Multiple p's in main",
			body: `
			<html><body>
				<p>Outside paragraph</p>
				<main>
					<p>Inside paragraph</p>
					<p>Inside paragraph 2</p>
				</main>
			</body></html>
			`,
			expected: "Inside paragraph",
			err:      false,
		},
		{
			name: "No p in main",
			body: `
			<html><body>
				<p>Outside paragraph</p>
				<main>
				</main>
			</body></html>
			`,
			expected: "Outside paragraph",
			err:      false,
		},
		{
			name: "No main",
			body: `
			<html><body>
				<p>Outside paragraph</p>
				<p>Outside paragraph 2</p>
			</body></html>
			`,
			expected: "Outside paragraph",
			err:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fPar, err := getFirstParagraphFromHTML(tc.body)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, fPar)
		})
	}
}
