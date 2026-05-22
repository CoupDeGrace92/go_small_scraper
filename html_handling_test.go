package main

import (
	"net/url"
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

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		htmlBody string
		expected []string
		err      bool
	}{
		{
			name:    "Base case",
			baseURL: "https://crawler-test.com",
			htmlBody: `
			<html><body>
				<a href="https://crawler-test.com"><span>RandomSite.com</span></a>
			</body></html>
			`,
			expected: []string{"https://crawler-test.com"},
			err:      false,
		},
		{
			name:    "Path appended",
			baseURL: "https://crawler-test.com",
			htmlBody: `
			<html><body>
				<a href="/path/coolstuff"><span>RandomSite.com</span></a>
			</body></html>
			`,
			expected: []string{"https://crawler-test.com/path/coolstuff"},
			err:      false,
		},
		{
			name:    "External link",
			baseURL: "https://crawler-test.com",
			htmlBody: `
			<html><body>
				<a href="https://google.com"><span>RandomSite.com</span></a>
			</body></html>
			`,
			expected: []string{"https://google.com"},

			err: false,
		},
		{
			name:    "Multiple links",
			baseURL: "https://crawler-test.com",
			htmlBody: `
			<html><body>
				<a href="https://crawler-test.com"><span>RandomSite.com</span></a>
				<a href="https://crawler-test.com/another_one"><span>RandomSite.com</span></a>
			</body></html>
			`,
			expected: []string{"https://crawler-test.com", "https://crawler-test.com/another_one"},
			err:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			base, err := url.Parse(tc.baseURL)
			require.NoError(t, err)
			urls, err := getURLsFromHTML(tc.htmlBody, base)
			if tc.err {
				require.Error(t, err)
				return
			}
			assert.ElementsMatch(t, tc.expected, urls)
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		htmlBody string
		expected []string
		err      bool
	}{
		{
			name:    "Base test",
			baseURL: "https://crawler-test.com",
			htmlBody: `
			<html><body>
				<img src="https://crawler-test.com/base_test.jpg" alt="Logo">
			</body></html>
			`,
			expected: []string{"https://crawler-test.com/base_test.jpg"},
			err:      false,
		},
		{
			name:    "path",
			baseURL: "https://crawler-test.com",
			htmlBody: `
			<html><body>
				<img src="/images/base_test.jpg" alt="Logo">
			</body></html>
			`,
			expected: []string{"https://crawler-test.com/images/base_test.jpg"},
			err:      false,
		},
		{
			name:    "External",
			baseURL: "https://crawler-test.com",
			htmlBody: `
			<html><body>
				<img src="https://image-site.com/base_test.jpg" alt="Logo">
			</body></html>
			`,
			expected: []string{"https://image-site.com/base_test.jpg"},
			err:      false,
		},
		{
			name:    "Multiple",
			baseURL: "https://crawler-test.com",
			htmlBody: `
			<html><body>
				<img src="https://crawler-test.com/base_test.jpg" alt="Logo">
				<img src="/another_one/base_test2.jpg" alt="Logo2">
			</body></html>
			`,
			expected: []string{"https://crawler-test.com/base_test.jpg", "https://crawler-test.com/another_one/base_test2.jpg"},
			err:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			base, err := url.Parse(tc.baseURL)
			require.NoError(t, err)
			urls, err := getImagesFromHTML(tc.htmlBody, base)
			if tc.err {
				require.Error(t, err)
				return
			}
			assert.ElementsMatch(t, tc.expected, urls)
		})
	}
}
