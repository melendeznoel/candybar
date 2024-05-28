package worldHealthOrg

import (
	"net/http"
)

// GenericResponse is the image compate response contract
type GenericResponse struct {
	ResponseError error
	Response      *http.Response
}

type Dataset struct {
	Label   string
	Display string
}

type Attribute struct {
	Label   string
	Display string
}

type Code struct {
	Label   string
	Display string
	Url     string
}

type Dimension struct {
	Label     string
	Display   string
	IsMeasure bool
	Code      []Code
}

type Dim struct {
	Category string
	Code     string
}

type Value struct {
	Display string
	Numeric string
	Low     string
	High    string
}

type Fact struct {
	Dateset         string
	Effectiive_Date string
	End_Date        string
	Pubished        bool
	Dim             []Dim
	Value
}

type InfantNutrition struct {
	Copyright string
	Dataset   []Dataset
	Attribute []Attribute
	Dimension []Dimension
	Fact      []Fact
}
