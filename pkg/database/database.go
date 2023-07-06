package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func GetDatabase(url string) (*sql.DB, error) {
	fmt.Println("Connecting to database...")
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	fmt.Println("Creating tables...")
	b, err := os.ReadFile("pkg/database/schema.sql")
	if err != nil {
		return nil, err
	}
	runSchema(db, string(b))
	return db, nil
}

func MustGetDatabase(url string) *sql.DB {
	db, err := GetDatabase(url)
	if err != nil {
		panic(err)
	}
	return db
}

func runSchema(db *sql.DB, schema string) {
	_, err := db.Exec(schema)
	if err != nil {
		panic(err)
	}
}
