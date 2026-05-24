package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	clArgs := os.Args[1:]
	if len(clArgs) < 1 {
		log.Printf("no website provided")
		os.Exit(1)
	} else if len(clArgs) > 3 {
		log.Printf("too many arguments provided")
		os.Exit(1)
	}
	fmt.Printf("starting crawl of: %v\n", clArgs[0])
	pages := make(map[string]PageData)
	urlObj, err := url.Parse(clArgs[0])
	if err != nil {
		fmt.Println("Error parsing url: ", err)
		os.Exit(1)
	}
	var buffSize int
	if len(clArgs) < 2 {
		buffSize = 5
	} else {
		b, err := strconv.Atoi(clArgs[1])
		if err != nil {
			fmt.Println("Could not convert buffsize arg to a number: ", err)
			fmt.Println("Defaulting max 5 active goroutines")
			buffSize = 5
		} else {
			buffSize = b
		}
	}
	var mPages int
	if len(clArgs) < 3 {
		mPages = 100
	} else {
		m, err := strconv.Atoi(clArgs[2])
		if err != nil {
			fmt.Println("Could not convert buffsize arg to a number: ", err)
			fmt.Println("Defaulting max 100 pages")
			m = 100
		} else {
			mPages = m
		}
	}
	cfg := Config{
		pages:      pages,
		baseURL:    urlObj,
		conControl: make(chan struct{}, buffSize),
		wg:         &sync.WaitGroup{},
		mu:         &sync.Mutex{},
		maxPages:   mPages,
	}
	cfg.wg.Add(1)
	cfg.crawlPage(clArgs[0])
	cfg.wg.Wait()
	for key, _ := range pages {
		fmt.Println(key)
	}
	writeJSONReport(cfg.pages, "report.json")
}
