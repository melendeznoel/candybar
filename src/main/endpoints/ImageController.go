package endpoints

import (
	"encoding/json"
	"net/http"

	ImageService "main/media"
)

// CompareImages returns stats on the image compare result
func CompareImages(w http.ResponseWriter, r *http.Request) {
	var figures []ImageService.Figure

	if r.Body == nil {
		http.Error(w, "Send Body", 400)
	}

	err := json.NewDecoder(r.Body).Decode(&figures)

	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	result := ImageService.Compare(figures)

	js, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
}
