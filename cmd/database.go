package main

import (
	"database/sql"
	"fmt"
)

func newDBDefaultSql() (*sql.DB, error) {
	// Create a PostgreSQL database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		DB_HOST,
		DB_USERNAME,
		DB_PASSWORD,
		DB_NAME,
		DB_PORT,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
