package server

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/raphael-foliveira/fiber-todo/pkg/common"
	"github.com/raphael-foliveira/fiber-todo/pkg/todo"
)

func StartServer(db *sql.DB) {
	fmt.Println("Starting server...")
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	startRoutes(app, db)
	app.Listen(":3000")
}

func startRoutes(app *fiber.App, db *sql.DB) {
	app.Get("/", common.StatusCheck)
	todoRoutes := app.Group("/todo")
	todo.GetTodoRoutes(todoRoutes, db)
}
