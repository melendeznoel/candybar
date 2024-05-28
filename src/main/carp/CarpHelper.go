package carp

import (
	"golang.org/x/net/html"
)

func scrapeHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val

			ok = true
		}
	}

	return
}
