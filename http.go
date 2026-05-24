package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Config struct {
	baseURL    *url.URL
	pages      map[string]PageData
	mu         *sync.Mutex
	conControl chan struct{}
	wg         *sync.WaitGroup
	maxPages   int
}

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

func (cfg *Config) crawlPage(rawCurrentURL string) {
	defer func() {
		<-cfg.conControl
		cfg.wg.Done()
	}()
	cfg.conControl <- struct{}{}

	if len(cfg.pages) >= cfg.maxPages {
		return
	}

	fmt.Printf("Crawling %s\n", rawCurrentURL)

	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("Error trying to parse %s into a url: %v\n", rawCurrentURL, err)
		return
	}
	if cfg.baseURL.Host != current.Host {
		return
	}

	normCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("Error normalizing url: %v\n", err)
		return
	}
	cfg.mu.Lock()
	_, exists := cfg.pages[normCurrent]
	var html string
	if !exists {
		html, err = getHTML(rawCurrentURL)
		if err != nil {
			log.Printf("Error fetching HTML: %v\n", err)
			return
		}
		pd, err := extractPageData(html, rawCurrentURL)
		if err != nil {
			log.Printf("Error extracting page data: %v\n", err)
			return
		}
		cfg.pages[normCurrent] = pd

	}
	cfg.mu.Unlock()
	//We divided this out so we could unlock slightly earlier
	if !exists {
		urls, err := getURLsFromHTML(html, cfg.baseURL)
		if err != nil {
			log.Printf("Error fetching urls from HTML: %v\n", err)
			return
		}
		for _, u := range urls {

			cfg.wg.Add(1)
			go func() {
				cfg.crawlPage(u)
			}()
		}
	}
}
