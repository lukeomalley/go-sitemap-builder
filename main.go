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
	maxDepth := flag.Int("max-depth", 5, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	for _, page := range pages {
		fmt.Println(page)
	}
}

type empty struct{}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]empty)

	var q map[string]empty
	nq := map[string]empty{
		urlStr: empty{},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]empty)
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = empty{}

			for _, link := range get(url) {
				nq[link] = empty{}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlString string) []string {
	resp, err := http.Get(urlString)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	base := baseURL(resp)
	return filterURLs(hrefsToLinks(resp.Body, base), withPrefix(base))
}

func baseURL(resp *http.Response) string {
	reqURL := resp.Request.URL
	b := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	return b.String()
}

func hrefsToLinks(r io.Reader, base string) []string {
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}

	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(
			l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}

	return hrefs
}

func filterURLs(links []string, keepFn func(string) bool) []string {
	var ret []string

	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
