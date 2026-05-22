package main

import (
	"log"
	"net/url"
	"strings"
)

func normalizeURL(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Printf("Error parsing url: %v\n", err)
		return "", err
	}

	u.Scheme = ""
	result := strings.TrimPrefix(u.String(), "//")
	result = strings.TrimSuffix(result, "/")
	return result, nil
}
