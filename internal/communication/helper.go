package communication

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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

// Logger will log a request info
func Logger(handler http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler.ServeHTTP(rw, r)

		log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
	})
}

// GetQueryParam returns parameter value
func GetQueryParam(paramKey string, urlVal *url.URL) (string, error) {
	if rawQuery, qe := url.ParseQuery(urlVal.RawQuery); qe != nil {
		return "", errors.New("Error when Raw Query")

	} else {
		if paramValue := rawQuery.Get(paramKey); paramValue != "" {
			return paramValue, nil

		}

		return "", errors.New(paramKey + " Parameter was not found")
	}
}

func GetParamsFromQuery(url *url.URL, key string) (result []string, found bool) {
	var (
		foundParams bool
		ids         []string
	)

	queryVals := url.Query()

	if len(queryVals[key]) > 0 {
		foundParams = true
	}

	for _, val := range queryVals[key] {
		ids = append(ids, val)
	}

	result = append(result, ids...)

	return result, foundParams
}

func RouteParam(r *http.Request, key string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}

	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	for i := 0; i < len(segments)-1; i++ {
		switch segments[i] {
		case "recipes", "ingredients":
			return segments[i+1]
		}
	}

	if len(segments) > 0 {
		return segments[len(segments)-1]
	}

	return ""
}
