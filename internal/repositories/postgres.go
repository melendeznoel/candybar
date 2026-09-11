package repositories

import (
	"database/sql"
	"fmt"
)

func All(db *sql.DB, tableName string) ([]interface{}, error) {
	var result []interface{}

	queryString := fmt.Sprintf("SELECT * FROM public.%s ORDER BY id ASC", tableName)

	rows, err := db.Query(queryString)

	if err != nil {
		return nil, fmt.Errorf("could not get all from %s: %s", tableName, err)
	}

	defer rows.Close()

	for rows.Next() {
		var item interface{}

		result = append(result, item)
	}

	return result, rows.Err()
}
