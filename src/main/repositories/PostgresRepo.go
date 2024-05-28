package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Db() (*sql.DB, error) {
	// todo: add to env
	// export POSTGRES_CONNECTION="user=[USER_NAME] password=[PASSWORD] dbname=[DATABASE_NAME] sslmode=disable

	datastoreName := os.Getenv("POSTGRES_CONNECTION")

	if db, err := sql.Open("postgres", datastoreName); err == nil {
		return db, nil
	} else {
		log.Printf("failed to get db instance: %v", err)

		return db, err
	}
}

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
