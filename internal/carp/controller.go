package carp

import (
	"encoding/json"
	"log"
	"net/http"
)

// Crawl endpoint
func Crawl(w http.ResponseWriter, r *http.Request) {
	// TODO: REPLACE
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		// todo: add middle ware to support this.
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	} else {
		urls, ok := r.URL.Query()["url"]

		if !ok || len(urls) < 1 {
			log.Println("url param is missing")

			w.WriteHeader(404)

			json.NewEncoder(w).Encode("urls values missing")

		} else {
			results := crawl(urls)

			json.NewEncoder(w).Encode(results)
		}
	}
}
