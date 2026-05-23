package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	clArgs := os.Args[1:]
	if len(clArgs) < 1 {
		log.Printf("no website provided")
		os.Exit(1)
	} else if len(clArgs) > 1 {
		log.Printf("too many arguments provided")
		os.Exit(1)
	}
	fmt.Printf("starting crawl of: %v\n", clArgs[0])
	pages := make(map[string]int)
	crawlPage(clArgs[0], clArgs[0], pages)
	for key, _ := range pages {
		fmt.Println(key)
	}
}
