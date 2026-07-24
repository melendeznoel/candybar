package food

import (
	"log"
	"time"
)

func getRecipeByID(id string) (Recipe, error) {
	var (
		result Recipe
		err    error
	)

	param := QueryParam{key: "id", value: id}

	result, err = queryRecipe(param)

	if err != nil {
		log.Printf("Failed to get Recipe: %v", err)

		return result, err
	}

	return result, nil
}

func getRecipesByID(ids []string) (Recipes, error) {
	var (
		params []QueryParam
	)

	for _, id := range ids {
		params = append(params, QueryParam{key: "id", value: id})
	}

	if len(params) == 0 {
		recipes, err := queryAllRecipes()

		if err != nil {
			log.Printf("Failed to get All Recipe: %v", err)

			return nil, err
		}

		return recipes, nil
	}

	recipes, err := queryRecipes(params)

	if err != nil {
		log.Printf("Failed to get Recipe: %v", err)

		return nil, err
	}

	return recipes, nil
}

func saveRecipe(recipe Recipe) (Recipe, error) {
	var newDateTime = time.Now()

	recipe.DateCreated = newDateTime
	recipe.DateLastModified = newDateTime

	// if recipe.Ingredients == nil {
	// 	recipe.Ingredients = json.RawMessage([]byte("{}"))
	// }

	var recipes Recipes

	recipes = append(recipes, recipe)

	result, err := insertRecipes(recipes)

	if err != nil {
		log.Printf("insert recipe : %v", err)
		return recipe, err
	}

	recipe.ID = result[0].ID

	return recipe, nil
}

func modifyRecipe(recipeID int64, recipe Recipe) (Recipe, error) {
	recipe.ID = recipeID

	recipe.DateLastModified = time.Now()

	// if recipe.Ingredients == nil {
	// 	recipe.Ingredients = json.RawMessage([]byte("{}"))
	// }

	updated, err := updateRecipe(recipe)

	if err != nil {
		log.Printf("update recipe : %v", err)
		return recipe, err
	}

	return updated, nil
}
