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
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func extractPageData(html, pageURL string) (PageData, error) {
	var out PageData
	urlObj, err := url.Parse(pageURL)
	heading, err := getHeadingFromHTML(html)
	if err != nil {
		return out, err
	}
	p, err := getFirstParagraphFromHTML(html)
	if err != nil {
		return out, err
	}
	imgs, err := getImagesFromHTML(html, urlObj)
	if err != nil {
		return out, err
	}
	urls, err := getURLsFromHTML(html, urlObj)
	if err != nil {
		return out, err
	}
	out = PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: p,
		OutgoingLinks:  urls,
		ImageURLs:      imgs,
	}
	return out, nil
}
