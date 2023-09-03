package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/raphael-foliveira/fiber-todo/docs"
	"github.com/raphael-foliveira/fiber-todo/pkg/common"
	"github.com/raphael-foliveira/fiber-todo/pkg/database"
	"github.com/raphael-foliveira/fiber-todo/pkg/todo"
)

// StartServer starts the server and adds the routes
func StartServer(db *database.Database) {
	fmt.Println("Starting server...")
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			err = c.Status(code).JSON(fiber.Map{
				"error":       err.Error(),
				"status_code": code,
			})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "internal server error",
				})
			}
			return nil
		},
	})
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())
	startRoutes(app, db)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatalf("%s", err)
	}
}

// startRoutes starts the routes for the application
func startRoutes(app *fiber.App, db *database.Database) {
	app.Get("/", common.StatusCheck)
	app.Get("/docs/*", swagger.HandlerDefault)
	apiRoutes := app.Group("/api")
	todoRoutes := apiRoutes.Group("/todos")
	repository := todo.NewTodoRepository(db)
	controller := todo.NewTodoController(repository)
	todo.GetTodoRoutes(todoRoutes, controller)
}
