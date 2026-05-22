package main

import (
	"strings"

	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	header := doc.Find("h1").Text()
	if header == "" {
		header = doc.Find("h2").Text()
	}

	return header, nil
}

func getFirstParagraphFromHTML(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	out := ""
	main := doc.Find("main")
	if main.Length() == 0 {
		out = doc.Find("p").First().Text()
	} else {
		p := main.Find("p")
		if p.Length() == 0 {
			out = doc.Find("p").First().Text()
		} else {
			out = p.First().Text()
		}
	}
	return out, nil
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	return []string{}, nil
}
