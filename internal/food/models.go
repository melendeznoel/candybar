package food

import (
	"time"
)

type Ingredient struct {
	ID int64

	Name        string
	Description string

	UnitOfWeight int64
	Quantity     int64
}

type Ingredients []Ingredient

type UnitOfWeight string

const (
	Cup          UnitOfWeight = "Cup"
	Dessertspoon UnitOfWeight = "Dessert Spoon"
	FluidOnce    UnitOfWeight = "Fluid Once"
	gallon       UnitOfWeight = "Gallon"
	Teaspoon     UnitOfWeight = "Teaspoon"
	Tablespoon   UnitOfWeight = "Tablespoon"
	Pint         UnitOfWeight = "Pint"
	Quart        UnitOfWeight = "Quart"
)

type Recipe struct {
	ID int64

	Name        string
	Description string

	Ingredients Ingredients

	// Workflows []Workflow

	// Ratings []Rating

	DateLastModified time.Time
	DateCreated      time.Time

	AverageRate int64
}

type Recipes []Recipe

type Workflow struct {
	Name        string
	Description string

	Index int

	Steps []Step
}

type Rating struct {
	ID   int
	Rate int

	Comment string

	DateCreated string
}

type Step struct {
	ID    int
	Index int

	Description string
	Notes       string

	Ingredients []Ingredient
}

type QueryParam struct {
	key   string
	value string
}

// Recall is a single openFDA food enforcement report
// (https://open.fda.gov/apis/food/enforcement/).
type Recall struct {
	RecallNumber            string `json:"recall_number"`
	ReasonForRecall         string `json:"reason_for_recall"`
	Status                  string `json:"status"`
	City                    string `json:"city"`
	State                   string `json:"state"`
	Country                 string `json:"country"`
	Classification          string `json:"classification"`
	ProductType             string `json:"product_type"`
	EventID                 string `json:"event_id"`
	RecallingFirm           string `json:"recalling_firm"`
	ReportDate              string `json:"report_date"`
	VoluntaryMandated       string `json:"voluntary_mandated"`
	InitialFirmNotification string `json:"initial_firm_notification"`
	DistributionPattern     string `json:"distribution_pattern"`
	ProductDescription      string `json:"product_description"`
	ProductQuantity         string `json:"product_quantity"`
	CodeInfo                string `json:"code_info"`
	TerminationDate         string `json:"termination_date"`
}

// RecallResultsMeta is the pagination/count block of an openFDA response.
type RecallResultsMeta struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// RecallMeta is the meta block of an openFDA response.
type RecallMeta struct {
	Disclaimer  string            `json:"disclaimer"`
	Terms       string            `json:"terms"`
	License     string            `json:"license"`
	LastUpdated string            `json:"last_updated"`
	Results     RecallResultsMeta `json:"results"`
}

// RecallResponse is the openFDA food enforcement API response envelope.
type RecallResponse struct {
	Meta    RecallMeta `json:"meta"`
	Results []Recall   `json:"results"`
}
