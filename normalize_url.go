package main

import (
	"errors"
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

	if u.Host == "" {
		return "", errors.New("missing host in URL")
	}

	u.Scheme = ""
	u.Host = strings.ToLower(u.Host)
	result := strings.TrimPrefix(u.String(), "//")
	result = strings.TrimSuffix(result, "/")
	return result, nil
}
