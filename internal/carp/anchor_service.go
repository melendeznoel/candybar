package carp

import (
	"strings"

	"golang.org/x/net/html"
)

func anchor(t html.Token, ch chan string) (ok bool, href string) {
	scraped, url := scrapeHref(t)

	if scraped {
		hasProto := strings.Index(url, "http") == 0

		if hasProto {
			ch <- url
		}

		ok = true
	}

	return
}
