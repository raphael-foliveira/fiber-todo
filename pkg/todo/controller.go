package todo

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raphael-foliveira/fiber-todo/pkg/common"
)

type TodoController struct {
	repository ITodoRepository
}

func NewTodoController(repository ITodoRepository) *TodoController {
	return &TodoController{repository: repository}
}

// @Create godoc
// @Summary Create a new To Do
// @Description Create a new To Do
// @Tags To Do
// @Accept json
// @Produce json
// @Param todo body CreateTodoDto true "To Do Create"
// @Success 201 {object} CreateResponse
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /todos [post]
func (tc *TodoController) Create(c *fiber.Ctx) error {
	todo, err := parseTodoFromBody(c)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, "bad request body")
	}
	createdTodo, err := tc.repository.Create(todo)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(fiber.StatusConflict, "todo already exists")
	}
	return c.Status(fiber.StatusCreated).JSON(createdTodo)
}

// @List godoc
// @Summary List To Dos
// @Description List To Dos
// @Tags To Do
// @Accept json
// @Produce json
// @Success 200 {array} Todo
// @Failure 500 {object} string "Internal Server Error"
// @Router /todos [get]
func (tc *TodoController) List(c *fiber.Ctx) error {
	todos, err := tc.repository.List()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(todos)
}

// @Retrieve godoc
// @Summary Retrieve a To Do
// @Description Retrieve a To Do
// @Tags To Do
// @Accept json
// @Produce json
// @Param id path int true "To Do ID"
// @Success 200 {object} Todo
// @Failure 404 {object} string "Not Found"
// @Failure 422 {object} string "Unprocessable Entity"
// @Router /todos/{id} [get]
func (tc *TodoController) Retrieve(c *fiber.Ctx) error {
	intId, err := common.ParseIdFromParams(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity)
	}
	todo, err := tc.repository.Retrieve(intId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(todo)
}

// @Update godoc
// @Summary Update a To Do
// @Description Update a To Do
// @Tags To Do
// @Accept json
// @Produce json
// @Param id path int true "To Do ID"
// @Param todo body UpdateTodoDto true "To Do Update"
// @Success 200 {object} Todo
// @Failure 404 {object} string "Not Found"
// @Failure 422 {object} string "Unprocessable Entity"
// @Failure 500 {object} string "Internal Server Error"
// @Router /todos/{id} [put]
func (tc *TodoController) Update(c *fiber.Ctx) (err error) {
	todo, err := parseTodoFromBody(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}
	todoId, err := common.ParseIdFromParams(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity)
	}
	uTodo, err := tc.repository.Update(Todo{Id: todoId, Title: todo.Title, Description: todo.Description, Completed: todo.Completed})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}
	if uTodo.Id == 0 {
		return fiber.NewError(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(uTodo)
}

// @Delete godoc
// @Summary Delete a To Do
// @Description Delete a To Do
// @Tags To Do
// @Accept json
// @Produce json
// @Param id path int true "To Do ID"
// @Success 204 "No Content"
// @Failure 422 {object} string "Unprocessable Entity"
// @Failure 500 {object} string "Internal Server Error"
// @Router /todos/{id} [delete]
func (tc *TodoController) Delete(c *fiber.Ctx) error {
	intId, err := common.ParseIdFromParams(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity)
	}
	affected, err := tc.repository.Delete(intId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}
	if affected == 0 {
		return fiber.NewError(fiber.StatusNotFound)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseTodoFromBody(c *fiber.Ctx) (CreateTodoDto, error) {
	var todo CreateTodoDto
	err := c.BodyParser(&todo)
	return todo, err
}
