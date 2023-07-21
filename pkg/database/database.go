package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	database "github.com/raphael-foliveira/fiber-todo/pkg/database/migrations"
)

func GetDatabase(url string) (*Database, error) {
	fmt.Println(url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	fmt.Println("connected to the database")
	return &Database{DB:db}, nil
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

func (db *Database) Migrate(schemaPath string) {
	fmt.Println("running migrations...")
	fmt.Println("schema:")
	fmt.Println(database.Schema)
	_, err := db.Exec(database.Schema)
	if err != nil {
		fmt.Println("error creating schema")
		panic(err)
	}

	fmt.Println("migrations:")
	for _, migration := range database.Migrations {
		fmt.Println(migration)
		_, err := db.Exec(migration)
		if err != nil {
			fmt.Println("error running migrations")
			panic(err)
		}
	}
	fmt.Println("created tables")
}
