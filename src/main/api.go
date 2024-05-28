package main

import (
	"log"
	"net/http"

	"main/endpoints"

	Carp "main/carp"
	IngredientController "main/food"
	RecipeController "main/food"
	HelperModels "main/helper"
	HelperService "main/helper"
	TwitterController "main/social"
	WorldHealthOrgController "main/worldHealthOrg"
)

// BuildRoutes returns a list of Route Objects
func BuildRoutes() HelperModels.Routes {
	return HelperService.Routes{
		{
			Name:            "Test",
			Method:          "GET",
			Pattern:         "/",
			HandlerFunction: endpoints.Hola,
		},
		{
			Name:            "CompareImages",
			Method:          "POST",
			Pattern:         "/image/compare",
			HandlerFunction: endpoints.CompareImages,
		},
		{
			Name:            "InfantNutrition",
			Method:          "GET",
			Pattern:         "/health/infantnutrition",
			HandlerFunction: WorldHealthOrgController.GetInfantNutrition,
		},
		{
			Name:            "TwitterHomeTimeline",
			Method:          "GET",
			Pattern:         "/tweets/hometimeline",
			HandlerFunction: TwitterController.GetHomeTimeline,
		},
		{
			Name:            "TwitterUserTimeline",
			Method:          "GET",
			Pattern:         "/tweets/usertimeline",
			HandlerFunction: TwitterController.GetUserTimeline,
		},
		{
			Name:            "Carp",
			Method:          "GET",
			Pattern:         "/carp/scrape",
			HandlerFunction: Carp.Crawl,
		},
		{
			Name:            "GetRecipes",
			Method:          "GET",
			Pattern:         "/food/recipes",
			HandlerFunction: RecipeController.GetRecipes,
		},
		{
			Name:            "PostRecipes",
			Method:          "POST",
			Pattern:         "/food/recipes",
			HandlerFunction: RecipeController.PostRecipes,
		},
		{
			Name:            "PutRecipes",
			Method:          "PUT",
			Pattern:         "/food/recipes/{id}",
			HandlerFunction: RecipeController.PutRecipes,
		},
		{
			Name:            "GetRecipeIngredients",
			Method:          "GET",
			Pattern:         "/food/recipes/{id}/ingredients",
			HandlerFunction: IngredientController.GetIngredients,
		},
		{
			Name:            "PostIngredients",
			Method:          "POST",
			Pattern:         "/food/ingredients",
			HandlerFunction: IngredientController.PostIngredients,
		},
		{
			Name:            "GetAllIngredients",
			Method:          "GET",
			Pattern:         "/food/ingredients",
			HandlerFunction: IngredientController.GetAllIngredients,
		},
	}
}

func main() {
	//appengine.Main()

	routes := BuildRoutes()

	router := HelperService.BuildRouter(routes)

	log.Fatal(http.ListenAndServe(":8080", router))
}

//todo: enable for app engine context
// func testhandle(w http.ResponseWriter, r *http.Request) {
// 	ctx := appengine.NewContext(r)

// 	googlelog.Infof(ctx, "got appengine context")
// }
