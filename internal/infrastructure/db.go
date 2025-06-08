package infrastructure

import (
	"database/sql"
	"log"
	"os"
)

func NewMySQL() *sql.DB {
	dsn := os.Getenv("DB_DSN")
	dbClient := os.Getenv("DB_CLIENT")

	db, err := sql.Open(dbClient, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
