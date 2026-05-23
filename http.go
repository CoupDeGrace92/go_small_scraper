package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", nil
	}

	req.Header.Set("User-Agent", "Simple-go-crawler v1.0")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()
	if resp.StatusCode > 400 {
		return "", fmt.Errorf("Status code indicates error: %v\n", resp.StatusCode)
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", errors.New("Content type not text/html")
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(htmlBytes), nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	base, err := url.Parse(rawBaseURL)
	if err != nil {
		log.Printf("Error trying to parse %s into a url: %v\n", rawBaseURL, err)
		return
	}
	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("Error trying to parse %s into a url: %v\n", rawCurrentURL, err)
		return
	}
	if base.Host != current.Host {
		return
	}

	normCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("Error normalizing url: %v\n", err)
		return
	}
	pages[normCurrent]++
	if pages[normCurrent] == 1 {
		html, err := getHTML(rawCurrentURL)
		if err != nil {
			log.Printf("Error fetching HTML: %v\n", err)
			return
		}
		b, err := url.Parse(rawBaseURL)
		if err != nil {
			log.Printf("Error trying to parse %s into a url object: %v\n", normCurrent, err)
			return
		}
		urls, err := getURLsFromHTML(html, b)
		if err != nil {
			log.Printf("Error fetching urls from HTML: %v\n", err)
			return
		}
		for _, u := range urls {
			crawlPage(rawBaseURL, u, pages)
		}
	}
}
