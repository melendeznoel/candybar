package helper

import (
	"net/http"
)

// Route object
type Route struct {
	Name            string
	Method          string
	Pattern         string
	HandlerFunction http.HandlerFunc
}

// Routes ...
type Routes []Route
