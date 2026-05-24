package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) {
	keys := make([]string, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pd []PageData
	for _, k := range keys {
		pd = append(pd, pages[k])
	}

	data, err := json.MarshalIndent(pd, "", "  ")
	if err != nil {
		log.Println("Error marshalling json: ", err)
	}

	os.WriteFile(filename, data, 0644)
}
