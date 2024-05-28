package endpoints

import (
	"encoding/json"
	"net/http"
)

// Hola returns an arrays of users
func Hola(w http.ResponseWriter, r *http.Request) {
	users := Users{User{FirstName: "mockfirstbane"}, User{FirstName: "secondmock", LastName: "mocklastname"}}

	json.NewEncoder(w).Encode(users)
}

// GetOnly validate the Request as a GET
func GetOnly(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			h(w, r)
			return
		}

		http.Error(w, "GET ONLY", http.StatusMethodNotAllowed)
	}
}

// PostOnly validate the Request as a POST
func PostOnly(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			h(w, r)
			return
		}

		http.Error(w, "POST ONLY", http.StatusMethodNotAllowed)
	}
}
