package worldHealthOrg

import (
	"encoding/json"
	"net/http"

	Helper "main/helper"
)

//GetInfantNutrition endpoint
func GetInfantNutrition(w http.ResponseWriter, r *http.Request) {
	// TODO: need a way to do this in one place.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if countryName, gqe := Helper.GetQueryParam("country", r.URL); gqe == nil {
		if data, e := FetchInfantNutrition(countryName); e == nil {
			json.NewEncoder(w).Encode(data)

		} else {
			w.WriteHeader(500)

			json.NewEncoder(w).Encode(e.Error())
		}

	} else {
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(gqe.Error())
	}
}
