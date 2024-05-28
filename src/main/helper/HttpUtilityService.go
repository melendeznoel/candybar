package helper

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// GetQueryParam returns parameter value
func GetQueryParam(paramKey string, urlVal *url.URL) (string, error) {
	if rawQuery, qe := url.ParseQuery(urlVal.RawQuery); qe != nil {
		return "", errors.New("Error when Raw Query")

	} else {
		if paramValue := rawQuery.Get(paramKey); paramValue != "" {
			return paramValue, nil

		}

		return "", errors.New(paramKey + " Parameter was not found")
	}
}

func GetParamsFromQuery(url *url.URL, key string) (result []string, found bool) {
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

func RouteParam(r *http.Request, key string) string {
	vars := mux.Vars(r)

	return vars[key]
}
