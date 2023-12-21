package todo

import (
	"errors"

	"github.com/raphael-foliveira/fiber-todo/pkg/database"
)

type ITodoRepository interface {
	Create(todo CreateTodoDto) (*Todo, error)
	List() ([]Todo, error)
	Retrieve(id int) (*Todo, error)
	Update(todo Todo) (*Todo, error)
	Delete(id int) (int64, error)
}

type TodoRepository struct {
	Db *database.Database
}

func NewTodoRepository(db *database.Database) *TodoRepository {
	return &TodoRepository{Db: db}
}

func (tr *TodoRepository) Create(todo CreateTodoDto) (*Todo, error) {
	tx, err := tr.Db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`
	INSERT INTO todo 
		(title, description, completed) 
	VALUES 
		($1, $2, $3) 
	RETURNING id, title, description, completed
	`)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(todo.Title, todo.Description, todo.Completed)
	var createdTodo Todo
	err = row.Scan(&createdTodo.Id, &createdTodo.Title, &createdTodo.Description, &createdTodo.Completed)
	return &createdTodo, err
}

func (tr *TodoRepository) List() ([]Todo, error) {
	rows, err := tr.Db.Query("SELECT id, title, description, completed FROM todo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := []Todo{}
	var todo Todo
	for rows.Next() {
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (tr *TodoRepository) Retrieve(id int) (*Todo, error) {
	row := tr.Db.QueryRow("SELECT id, title, description, completed FROM todo WHERE id = $1", id)
	var todo Todo
	err := row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (tr *TodoRepository) Update(todo Todo) (*Todo, error) {
	result, err := tr.Db.Exec("UPDATE todo SET title = $1, description = $2, completed = $3 WHERE id = $4",
		todo.Title, todo.Description, todo.Completed, todo.Id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("todo not found")
	}
	return &todo, nil
}

func (tr *TodoRepository) Delete(id int) (int64, error) {
	result, err := tr.Db.Exec("DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if affectedRows == 0 {
		return 0, errors.New("todo not found")
	}
	return affectedRows, nil
}
