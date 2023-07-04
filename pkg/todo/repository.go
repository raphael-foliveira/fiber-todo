package todo

import "database/sql"

type TodoRepository struct {
	Db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{Db: db}
}

func (tr *TodoRepository) Create(todo Todo) (int, error) {
	row := tr.Db.QueryRow("INSERT INTO todo (title, description, completed) VALUES ($1, $2, $3) RETURNING id",
		todo.Title, todo.Description, todo.Completed)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (tr *TodoRepository) List() ([]Todo, error) {
	rows, err := tr.Db.Query("SELECT id, title, description, completed FROM todo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := []Todo{}
	for rows.Next() {
		todo := Todo{}
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (tr *TodoRepository) Retrieve(id int) (Todo, error) {
	row := tr.Db.QueryRow("SELECT id, title, description, completed FROM todo WHERE id = $1", id)
	var todo Todo
	err := row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (tr *TodoRepository) Update(todo Todo) (Todo, error) {
	result, err := tr.Db.Exec("UPDATE todo SET title = $1, description = $2, completed = $3 WHERE id = $4",
		todo.Title, todo.Description, todo.Completed, todo.Id)
	if err != nil {
		return Todo{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Todo{}, err
	}
	if rowsAffected == 0 {
		return Todo{}, nil
	}
	return todo, nil
}

func (tr *TodoRepository) Delete(id int) error {
	_, err := tr.Db.Exec("DELETE FROM todo WHERE id = $1", id)
	return err
}
