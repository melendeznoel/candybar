package worldhealthorg

import (
	"net/http"
)

// GenericResponse is the image compate response contract
type GenericResponse struct {
	ResponseError error
	Response      *http.Response
}

// GhoObservation is a single data point as returned by the GHO OData API
// (https://www.who.int/data/gho/info/gho-odata-api), which replaced the
// retired Athena API.
type GhoObservation struct {
	ID                 int     `json:"Id"`
	IndicatorCode      string  `json:"IndicatorCode"`
	SpatialDimType     string  `json:"SpatialDimType"`
	SpatialDim         string  `json:"SpatialDim"`
	ParentLocationCode string  `json:"ParentLocationCode"`
	ParentLocation     string  `json:"ParentLocation"`
	TimeDimType        string  `json:"TimeDimType"`
	TimeDim            int     `json:"TimeDim"`
	Dim1Type           string  `json:"Dim1Type"`
	Dim1               string  `json:"Dim1"`
	NumericValue       float64 `json:"NumericValue"`
	Low                float64 `json:"Low"`
	High               float64 `json:"High"`
	Value              string  `json:"Value"`
	Comments           string  `json:"Comments"`
	Date               string  `json:"Date"`
	TimeDimensionValue string  `json:"TimeDimensionValue"`
	TimeDimensionBegin string  `json:"TimeDimensionBegin"`
	TimeDimensionEnd   string  `json:"TimeDimensionEnd"`
}

// InfantNutrition is the GHO OData API response envelope for the
// WHOSIS_000006 indicator (infants exclusively breastfed for the first six
// months of life).
type InfantNutrition struct {
	Context string           `json:"@odata.context"`
	Value   []GhoObservation `json:"value"`
}

// InfantDeaths is the GHO OData API response envelope for the CM_02
// indicator (number of infant deaths).
type InfantDeaths struct {
	Context string           `json:"@odata.context"`
	Value   []GhoObservation `json:"value"`
}
