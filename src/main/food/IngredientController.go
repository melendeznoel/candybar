package food

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	HttpUtilityService "main/helper"
)

func GetIngredients(w http.ResponseWriter, r *http.Request) {
	var (
		id    int64
		idStr string

		err error
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	} else {
		idStr = HttpUtilityService.RouteParam(r, "id")

		if id, err = strconv.ParseInt(idStr, 10, 64); err != nil {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(err)
			return
		}

		if id == 0 {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode("id should be non-zero")
		}

		if err, response := getIngredientsByRecipeID(id); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(response)
		}
	}
}

func GetAllIngredients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	if err, response := getIngredientsByRecipeID(0); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(response)
	}
}

func PostIngredients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	} else {
		json.NewEncoder(w).Encode(Ingredient{Name: "mock"})
	}
}

func ingredientIdsFromQuery(url *url.URL, key string) (result []string, found bool) {
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
