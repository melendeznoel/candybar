package worldHealthOrg

import (
	"net/http"
)

// Get is generic http get service
func Get(url string) GenericResponse {
	r, e := http.Get(url)

	return GenericResponse{Response: r, ResponseError: e}
}
