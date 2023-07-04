package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/raphael-foliveira/fiber-todo/pkg/database"
	"github.com/raphael-foliveira/fiber-todo/pkg/server"
)

func main() {
	godotenv.Load()
	db := database.MustGetDatabase(os.Getenv("DATABASE_URL"), "pkg/database/schema.sql")
	defer db.Close()
	server.StartServer(db)
}
