package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func GetDatabase(url string) (*sql.DB, error) {
	return sql.Open("postgres", url)
}

func MustGetDatabase(url string) *sql.DB {
	db, err := GetDatabase(url)
	if err != nil {
		panic(err)
	}
	return db
}
