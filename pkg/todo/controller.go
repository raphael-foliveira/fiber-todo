package todo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raphael-foliveira/fiber-todo/pkg/common"
)

type TodoController struct {
	repository ITodoRepository
}

func NewTodoController(repository ITodoRepository) *TodoController {
	return &TodoController{repository: repository}
}

func (tc *TodoController) Create(c *fiber.Ctx) error {
	todo, err := parseTodoFromBody(c)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	id, err := tc.repository.Create(todo)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (tc *TodoController) List(c *fiber.Ctx) error {
	todos, err := tc.repository.List()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(todos)
}

func (tc *TodoController) Retrieve(c *fiber.Ctx) error {
	intId, err := common.ParseIdFromParams(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	todo, err := tc.repository.Retrieve(intId)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(todo)
}

func (tc *TodoController) Update(c *fiber.Ctx) (err error) {
	todo, err := parseTodoFromBody(c)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	todo.Id, err = common.ParseIdFromParams(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	uTodo, err := tc.repository.Update(todo)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if uTodo.Id == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(uTodo)
}

func (tc *TodoController) Delete(c *fiber.Ctx) error {
	intId, err := common.ParseIdFromParams(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	err = tc.repository.Delete(intId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseTodoFromBody(c *fiber.Ctx) (Todo, error) {
	var todo Todo
	if err := c.BodyParser(&todo); err != nil {
		return todo, err
	}
	return todo, nil
}
