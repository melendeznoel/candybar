package worldHealthOrg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// FetchInfantNutrition returns Infant Nutrition
func FetchInfantNutrition(country string) (*InfantNutrition, error) {
	var in = new(InfantNutrition)

	response := Get("http://apps.who.int/gho/athena/data/GHO/WHS_PBR,WHOSIS_000006.json?filter=COUNTRY:" + country + ";REGION:*")

	if response.ResponseError != nil {
		return nil, errors.New("error on getting Infant Nutrition")
	}

	body, raErr := ioutil.ReadAll(response.Response.Body)

	if raErr != nil {
		return nil, raErr
	}

	umErr := json.Unmarshal([]byte(body), &in)

	if umErr != nil {
		fmt.Println("HS an Err:", umErr)
	}

	return in, nil
}
