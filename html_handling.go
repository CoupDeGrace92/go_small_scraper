package main

import (
	"fmt"
	"net/url"
	"strings"

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
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	out := []string{}
	doc.Find("a").Each(func(index int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists {
			//URL cleaning here
			ref, err := url.Parse(href)
			if err != nil {
				fmt.Printf("Error converting href into a url object: %v\n", err)
				return //in goquery - this will just move onto the next element
			}
			href = baseURL.ResolveReference(ref).String()
			out = append(out, href)
		}
	})
	return out, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	out := []string{}
	doc.Find("img").Each(func(index int, element *goquery.Selection) {
		src, exists := element.Attr("src")
		if exists {
			source, err := url.Parse(src)
			if err != nil {
				fmt.Printf("Error converting src into a url object: %v\n", err)
				return
			}
			src = baseURL.ResolveReference(source).String()
			out = append(out, src)
		}
	})

	return out, nil
}

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	return PageData{}
}
