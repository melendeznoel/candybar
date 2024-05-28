package food

func getIngredientsByRecipeID(recipeID int64) (Ingredients, error) {
	var (
		result Ingredients
		err    error
	)

	if recipeID == 0 {
		result, err = queryAllIngredients(nil)
	} else {
		result, err = queryIngredients(recipeID, nil)
	}

	return result, err
}
