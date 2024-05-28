package carp

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func crawl(urls []string) []string {
	// urlsFound := make(map[string]bool)
	urlsFound := []string{}

	// channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// scrape
	for _, url := range urls {
		go scrape(url, chUrls, chFinished)
	}

	for c := 0; c < len(urls); {
		select {
		case url := <-chUrls:
			// urlsFound[url] = true

			urlsFound = append(urlsFound, url)

		case <-chFinished:
			c++
		}
	}

	close(chUrls)

	return urlsFound
}

func scrape(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("Failed scrape \"" + url + "\"")

		return
	}

	b := resp.Body

	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:

			return

		case tt == html.StartTagToken:
			t := z.Token()

			// TODO: REMOVE
			// isAnchor := t.Data == "a"

			if !isAnchor(t) {
				continue
			}

			ok, url := scrapeHref(t)

			if !ok {
				continue
			}

			hasProto := strings.Index(url, "http") == 0

			if hasProto {
				ch <- url
			}
		}
	}
}

func isAnchor(t html.Token) bool {
	return t.Data == "a"
}
