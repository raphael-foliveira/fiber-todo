package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	_ "github.com/raphael-foliveira/fiber-todo/docs"
	"github.com/raphael-foliveira/fiber-todo/pkg/common"
	"github.com/raphael-foliveira/fiber-todo/pkg/todo"
)

// StartServer starts the server and adds the routes
func StartServer(db *sql.DB) {
	fmt.Println("Starting server...")
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	startRoutes(app, db)
	err := app.Listen(":3000")
	if err != nil {
		log.Fatalf("%s", err)
	}
}

// startRoutes starts the routes for the application
func startRoutes(app *fiber.App, db *sql.DB) {
	app.Get("/", common.StatusCheck)
	app.Get("/docs/*", swagger.HandlerDefault)
	apiRoutes := app.Group("/api")
	todoRoutes := apiRoutes.Group("/todos")
	todo.GetTodoRoutes(todoRoutes, db)
}
