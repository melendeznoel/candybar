package food

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// openFDAFoodEnforcementURL is the openFDA Food Enforcement (recall) endpoint
// (https://open.fda.gov/apis/food/enforcement/).
const openFDAFoodEnforcementURL = "https://api.fda.gov/food/enforcement.json"

// FetchRecallsByProductDescription returns openFDA food enforcement records
// whose product_description matches the given search term.
func FetchRecallsByProductDescription(productDescription string) (*RecallResponse, error) {
	query := url.Values{}
	query.Set("search", fmt.Sprintf("product_description:%s", productDescription))
	query.Set("limit", "100")
	query.Set("skip", "0")

	reqURL := openFDAFoodEnforcementURL + "?" + query.Encode()

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// openFDA responds 404 (with an error envelope) when no records match a
	// search, rather than 200 with an empty results array.
	if resp.StatusCode == http.StatusNotFound {
		return &RecallResponse{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openFDA returned status %d: %s", resp.StatusCode, string(body))
	}

	var result RecallResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// HasRecalls reports whether any openFDA recall records exist for the given
// product description.
func HasRecalls(productDescription string) (bool, error) {
	result, err := FetchRecallsByProductDescription(productDescription)
	if err != nil {
		return false, err
	}

	return len(result.Results) > 0, nil
}
