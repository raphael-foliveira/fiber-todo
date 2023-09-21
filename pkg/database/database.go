package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/raphael-foliveira/fiber-todo/pkg/database/queries"
)

func GetDatabase(url string) (*Database, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}

func MustGetDatabase(url string) *Database {
	db, err := GetDatabase(url)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	return db
}

type Database struct {
	*sql.DB
}

func (db *Database) CreateSchema() {
	_, err := db.Exec(queries.Schema)
	if err != nil {
		fmt.Println("error creating schema")
		panic(err)
	}
}
