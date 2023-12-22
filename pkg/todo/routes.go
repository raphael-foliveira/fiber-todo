package todo

import (
	"github.com/gofiber/fiber/v2"
)

func GetTodoRoutes(router fiber.Router, controller *TodoController) fiber.Router {
	router.Post("/", controller.Create)
	router.Get("/", controller.List)
	router.Get("/:id", controller.Retrieve)
	router.Put("/:id", controller.Update)
	router.Delete("/:id", controller.Delete)
	return router
}
