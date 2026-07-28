package worldhealthorg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// ghoOdataBaseURL is the GHO OData API (https://www.who.int/data/gho/info/gho-odata-api),
// which replaced the retired GHO Minerva interface and Athena API.
const ghoOdataBaseURL = "https://ghoapi.azureedge.net/api"

// FetchInfantNutrition returns Infant Nutrition
func FetchInfantNutrition(country string) (*InfantNutrition, error) {
	var infantNutrition = new(InfantNutrition)

	query := url.Values{}
	query.Set("$filter", fmt.Sprintf("SpatialDim eq '%s'", country))

	// OData filter expressions use spaces and single quotes that survive
	// url.Values encoding fine, but "+" for spaces reads oddly in $filter;
	// %20 matches the OData API's documented examples.
	reqURL := ghoOdataBaseURL + "/WHOSIS_000006?" + strings.ReplaceAll(query.Encode(), "+", "%20")

	response := Get(reqURL)

	if response.ResponseError != nil {
		return nil, errors.New("error on getting Infant Nutrition")
	}

	defer response.Response.Body.Close()

	body, raErr := io.ReadAll(response.Response.Body)

	if raErr != nil {
		return nil, raErr
	}

	if umErr := json.Unmarshal(body, &infantNutrition); umErr != nil {
		return nil, umErr
	}

	return infantNutrition, nil
}

// FetchInfantDeaths returns the number of infant deaths (indicator CM_02)
func FetchInfantDeaths(country string) (*InfantDeaths, error) {
	var infantDeaths = new(InfantDeaths)

	query := url.Values{}
	query.Set("$filter", fmt.Sprintf("SpatialDim eq '%s'", country))

	// OData filter expressions use spaces and single quotes that survive
	// url.Values encoding fine, but "+" for spaces reads oddly in $filter;
	// %20 matches the OData API's documented examples.
	reqURL := ghoOdataBaseURL + "/CM_02?" + strings.ReplaceAll(query.Encode(), "+", "%20")

	response := Get(reqURL)

	if response.ResponseError != nil {
		return nil, errors.New("error on getting Infant Deaths")
	}

	defer response.Response.Body.Close()

	body, raErr := io.ReadAll(response.Response.Body)

	if raErr != nil {
		return nil, raErr
	}

	if umErr := json.Unmarshal(body, &infantDeaths); umErr != nil {
		return nil, umErr
	}

	return infantDeaths, nil
}
