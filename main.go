package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/raphael-foliveira/fiber-todo/pkg/database"
	"github.com/raphael-foliveira/fiber-todo/pkg/server"
)

// @title           Fiber To Do API
// @version         1.0
// @description     A To Do app built with the Fiber framework

// @contact.name   Raphael Oliveira
// @contact.url    https://github.com/raphael-foliveira
// @BasePath /api
func main() {
	godotenv.Load()
	db := database.MustGetDatabase(os.Getenv("DATABASE_URL"))
	db.Migrate()
	defer db.Close()
	server.StartServer(db.DB)
}
