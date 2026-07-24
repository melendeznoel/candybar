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
