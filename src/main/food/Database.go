package food

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" //required by sql package
)

func db() (*sql.DB, error) {
	// todo: add to env
	// export POSTGRES_CONNECTION="user=[USER_NAME] password=[PASSWORD] dbname=[DATABASE_NAME] sslmode=disable

	datastoreName := os.Getenv("POSTGRES_CONNECTION")

	db, err := sql.Open("postgres", datastoreName)

	if err != nil {
		log.Printf("failed to get db instance: %v", err)

		return db, nil
	}

	return db, nil
}
