package food

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	HttpUtilityService "main/helper"
)

func GetRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	ids, found := recipeIDsFromQuery(r.URL, "id")

	if found == false {
		if recipes, err := getRecipesByID(nil); err != nil {
			w.WriteHeader(500)

			json.NewEncoder(w).Encode(err)

		} else {
			json.NewEncoder(w).Encode(recipes)
		}
		return
	}

	if len(ids) == 1 {
		if recipes, err := getRecipeByID(ids[0]); err != nil {
			w.WriteHeader(500)

			json.NewEncoder(w).Encode(err.Error())

		} else {
			json.NewEncoder(w).Encode(recipes)
		}
		return
	}

	if recipes, err := getRecipesByID(ids); err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(err)

	} else {
		json.NewEncoder(w).Encode(recipes)
	}
}

func PostRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var (
		recipe Recipe
	)

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	if r.Body == nil {
		err := errors.New("Body in Request missing")

		w.WriteHeader(400)

		json.NewEncoder(w).Encode(err)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&recipe)

	if err != nil {
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(err)
		return
	}

	data, err := saveRecipe(recipe)

	if err != nil {
		w.WriteHeader(422)

		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func PutRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var (
		recipe Recipe
		id     int64
	)

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		return
	}

	recipeID := HttpUtilityService.RouteParam(r, "id")

	if intID, err := strconv.ParseInt(recipeID, 10, 64); err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(err)
		return
	} else {
		//todo: replace this
		id = intID
	}

	if r.Body == nil {
		err := errors.New("Body in Request missing")

		w.WriteHeader(400)

		json.NewEncoder(w).Encode(err)
	}

	err := json.NewDecoder(r.Body).Decode(&recipe)

	if err != nil {
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(err)
		return
	}

	if _, err := modifyRecipe(id, recipe); err != nil {
		w.WriteHeader(422)

		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(recipe)
}

func recipeIDsFromQuery(url *url.URL, key string) (result []string, found bool) {
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
