package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func GetDatabase(url, schemaPath string) (*sql.DB, error) {
	fmt.Println("Connecting to database...")
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	fmt.Println("Creating tables...")
	b, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}
	runSchema(db, string(b))
	return db, nil
}

func MustGetDatabase(url, schemaPath string) *sql.DB {
	db, err := GetDatabase(url, schemaPath)
	if err != nil {
		fmt.Println(err.Error())
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
