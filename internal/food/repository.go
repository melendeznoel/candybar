package food

import (
	"errors"
	"fmt"
	"strconv"

	pq "github.com/lib/pq"
)

var tableName = "recipe"

func queryRecipe(queryParam QueryParam) (Recipe, error) {
	var (
		result Recipe
		found  bool
	)

	db, err := db()

	if err != nil {
		return result, err
	}

	sqlStatement := "SELECT * FROM public." + tableName + " WHERE " + queryParam.key + "=" + queryParam.value + " ORDER BY id ASC"

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return result, fmt.Errorf("could not get all from %s: %s", tableName, err)
	}

	defer rows.Close()

	for rows.Next() {
		item := Recipe{}

		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.DateLastModified, &item.DateCreated, &item.AverageRate, &item.Ingredients); err != nil {
			return result, fmt.Errorf("Could not get name/description out of row: %v", err)
		}

		result = item

		found = true
	}

	if !found {
		var noRowsErr = errors.New("Recipe not found")
		return result, noRowsErr
	}

	return result, nil
}

func queryRecipes(queryParam []QueryParam) (Recipes, error) {
	// todo: need to create batch get
	var (
		result Recipes
	)

	db, err := db()

	if err != nil {
		return nil, err
	}

	sqlStatement := `SELECT * FROM public.%s ORDER BY id ASC`

	rows, err := db.Query(sqlStatement, tableName, queryParam[0].key, queryParam[0].value)

	if err != nil {
		return nil, fmt.Errorf("could not get all from %s: %s", tableName, err)
	}

	defer rows.Close()

	for rows.Next() {
		var item Recipe

		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.DateLastModified, &item.DateCreated, &item.AverageRate, &item.Ingredients); err != nil {
			return nil, fmt.Errorf("Could not get name/description out of row: %v", err)
		}

		result = append(result, item)
	}

	return result, rows.Err()
}

func queryAllRecipes() (Recipes, error) {
	//todo: need to pss in query params
	var (
		result Recipes
	)

	db, err := db()

	if err != nil {
		return nil, err
	}

	sqlStatement := "SELECT * FROM public." + tableName + " ORDER BY id ASC"

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, fmt.Errorf("could not get all from %s: %s", tableName, err)
	}

	defer rows.Close()

	for rows.Next() {
		var item Recipe

		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.DateLastModified, &item.DateCreated, &item.AverageRate, &item.Ingredients); err != nil {
			return nil, fmt.Errorf("Could not get name/description out of row: %v", err)
		}

		result = append(result, item)
	}

	return result, rows.Err()
}

func insertRecipes(recipes Recipes) (Recipes, error) {
	db, err := db()

	if err != nil {
		return recipes, err
	}

	txn, err := db.Begin()

	if err != nil {
		return recipes, err
	}

	stmt, err := txn.Prepare(pq.CopyIn("recipe", "name", "description", "averageRate", "dateLastModified", "dateCreated"))

	if err != nil {
		return recipes, err
	}

	for _, recipe := range recipes {
		result, err := stmt.Exec(recipe.Name, recipe.Description, int64(recipe.AverageRate), recipe.DateCreated, recipe.DateLastModified)

		if err == nil {
			//todo: err handling
			id, _ := result.LastInsertId()

			recipe.ID = id
		}
	}

	_, err = stmt.Exec()

	if err != nil {
		return recipes, err
	}

	err = stmt.Close()

	if err != nil {
		return recipes, err
	}

	err = txn.Commit()

	if err != nil {
		return recipes, err
	}

	return recipes, nil
}

func updateRecipe(recipe Recipe) (Recipe, error) {
	db, err := db()

	if err != nil {
		return recipe, err
	}

	updateStatement := `
	UPDATE recipe
	SET name = $2, description = $3, ingredients = $4, "dateLastModified" = $5, "averageRate" = $6
	WHERE ID = $1;`

	_, err = db.Exec(updateStatement, recipe.ID, recipe.Name, recipe.Description, recipe.Ingredients, recipe.DateLastModified, recipe.AverageRate)

	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

func queryIngredients(recipeID int64, queryParams []QueryParam) (Ingredients, error) {
	var (
		err    error
		recipe Recipe
	)

	result := Ingredients{}

	recipeQueryParam := QueryParam{key: "ID", value: strconv.Itoa(int(recipeID))}

	recipe, err = queryRecipe(recipeQueryParam)

	if err != nil {
		return result, err
	}

	result = recipe.Ingredients

	return result, err
}

func queryAllIngredients(queryParams []QueryParam) (Ingredients, error) {
	var (
		err     error
		recipes Recipes
	)

	result := Ingredients{}

	recipes, err = queryAllRecipes()

	if err != nil {
		return result, err
	}

	for _, r := range recipes {
		result = append(result, r.Ingredients...)
	}

	//todo: filter using the queryParams

	return result, err
}
