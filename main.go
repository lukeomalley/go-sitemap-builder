package main

import (
	"flag"
	"fmt"
	"gophercises/html-link-parser/link"
	"io"
	"net/http"
	"net/url"
	"strings"
)

/*
	1. GET the webpage
	2. parse all the links on the page
	3. build proper urls with our links
	4. filter out links to a different domain
	5. find all the pages (BFS)
	6. Print out XML
*/

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	flag.Parse()

	fmt.Println(*urlFlag)

	// GET the webpage
	resp, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	base := baseURL(resp)
	pages := hrefsToPages(resp.Body, base)
	for _, h := range pages {
		fmt.Println(h)
	}
}

func baseURL(resp *http.Response) string {
	reqURL := resp.Request.URL
	b := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	return b.String()
}

func hrefsToPages(r io.Reader, base string) []string {
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}

	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}

	return hrefs
}
