package social

import (
	"encoding/json"
	"net/http"
)

// GetHomeTimeline endpoint
func GetHomeTimeline(w http.ResponseWriter, r *http.Request) {
	// TODO: need a way to do this in one place.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	} else {
		if tweets, e := FetchHomeTimeline(); e == nil {
			json.NewEncoder(w).Encode(tweets)

		} else {
			w.WriteHeader(500)

			json.NewEncoder(w).Encode(e.Error())
		}
	}
}

// GetUserTimeline endpoint
func GetUserTimeline(w http.ResponseWriter, r *http.Request) {
	// TODO: need a way to do this globally.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	} else {
		if tweets, e := FetchUserTimeline(); e == nil {
			json.NewEncoder(w).Encode(tweets)

		} else {
			w.WriteHeader(500)

			json.NewEncoder(w).Encode(e.Error())
		}
	}
}
