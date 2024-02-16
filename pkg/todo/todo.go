package todo

import "github.com/raphael-foliveira/fiber-todo/pkg/database"

type TodoModule struct {
	Repository ITodoRepository
	Controller *TodoController
}

func New(db *database.Database) *TodoModule {
	repository := NewTodoRepository(db)
	controller := NewTodoController(repository)
	return &TodoModule{
		Repository: repository,
		Controller: controller,
	}
}
