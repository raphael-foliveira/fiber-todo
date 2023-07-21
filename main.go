package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/raphael-foliveira/fiber-todo/pkg/common"
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
	args := os.Args[1:]
	godotenv.Load()
	db := database.MustGetDatabase(os.Getenv("DATABASE_URL"))
	if (common.Contains(args, "migrate")) {
		db.Migrate("./pkg/database/schema.sql")
		return
	}
	defer db.Close()
	server.StartServer(db.DB)
}
