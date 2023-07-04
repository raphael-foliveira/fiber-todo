package todo

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func GetTodoRoutes(group fiber.Router, db *sql.DB) fiber.Router {
	repository := NewTodoRepository(db)
	controller := NewTodoController(repository)
	group.Post("/", controller.Create)
	group.Get("/", controller.List)
	group.Get("/:id", controller.Retrieve)
	group.Put("/:id", controller.Update)
	group.Delete("/:id", controller.Delete)
	return group
}
