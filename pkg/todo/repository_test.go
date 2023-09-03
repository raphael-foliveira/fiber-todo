package todo

import (
	"testing"

	"github.com/raphael-foliveira/fiber-todo/pkg/common"
	"github.com/raphael-foliveira/fiber-todo/pkg/database"
	"github.com/raphael-foliveira/fiber-todo/pkg/database/queries"
)

var config = common.ReadTestCfg()

var repository *TodoRepository

func repositoryTestsSetup() {
	db := database.MustGetDatabase(config.Database.Url)
	db.Migrate()
	repository = NewTodoRepository(db)
}

func repositoryTestsTeardown() {
	repository.Db.Exec(queries.RecreateSchema)
	repository.Db.Close()
}

func TestRepositoryCreate(t *testing.T) {
	repositoryTestsSetup()
	defer repositoryTestsTeardown()
	todo, err := repository.Create(CreateTodoDto{
		Title:       "Test",
		Description: "Test",
		Completed:   false,
	})
	if err != nil {
		t.Errorf("Error creating todo: %s", err)
	}
	if todo.Title != "Test" {
		t.Errorf("Expected title to be 'Test', got '%s'", todo.Title)
	}
	if todo.Description != "Test" {
		t.Errorf("Expected description to be 'Test', got '%s'", todo.Description)
	}
	if todo.Completed != false {
		t.Errorf("Expected completed to be false, got '%t'", todo.Completed)
	}
}

func TestRepositoryList(t *testing.T) {
	repositoryTestsSetup()
	defer repositoryTestsTeardown()
	repository.Db.Exec(queries.InsertTodoFixtures)
	todos, err := repository.List()
	if err != nil {
		t.Errorf("Error listing todos: %s", err)
	}
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}
}

func TestRepositoryRetrieve(t *testing.T) {
	t.Run("should return the todo with the given id", func(t *testing.T) {
		repositoryTestsSetup()
		defer repositoryTestsTeardown()
		repository.Db.Exec(queries.InsertTodoFixtures)
		todo, err := repository.Retrieve(1)
		if err != nil {
			t.Errorf("Error retrieving todo: %s", err)
		}
		if todo.Id != 1 {
			t.Errorf("Expected id to be 1, got %d", todo.Id)
		}
	})

	t.Run("should return an error when given an id that doesn't exist", func(t *testing.T) {
		repositoryTestsSetup()
		defer repositoryTestsTeardown()
		_, err := repository.Retrieve(1)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestRepositoryUpdate(t *testing.T) {
	t.Run("should update the todo with the given id", func(t *testing.T) {
		repositoryTestsSetup()
		defer repositoryTestsTeardown()
		repository.Db.Exec(queries.InsertTodoFixtures)
		todo, err := repository.Update(Todo{
			Id:          1,
			Title:       "Updated",
			Description: "Updated",
			Completed:   true,
		})
		if err != nil {
			t.Errorf("Error updating todo: %s", err)
		}
		if todo.Title != "Updated" {
			t.Errorf("Expected title to be 'Updated', got '%s'", todo.Title)
		}
		if todo.Description != "Updated" {
			t.Errorf("Expected description to be 'Updated', got '%s'", todo.Description)
		}
		if todo.Completed != true {
			t.Errorf("Expected completed to be true, got '%t'", todo.Completed)
		}
	})

	t.Run("should return an error when given an id that doesn't exist", func(t *testing.T) {
		repositoryTestsSetup()
		defer repositoryTestsTeardown()
		repository.Db.Exec(queries.ClearTodoTable)
		_, err := repository.Update(Todo{
			Id:          9999,
			Title:       "Updated",
			Description: "Updated",
			Completed:   true,
		})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestRepositoryDelete(t *testing.T) {
	t.Run("should delete a todo with the given id", func(t *testing.T) {
		repositoryTestsSetup()
		defer repositoryTestsTeardown()
		repository.Db.Exec(queries.InsertTodoFixtures)
		rowsAffected, err := repository.Delete(1)
		if err != nil {
			t.Errorf("Error deleting todo: %s", err)
		}
		if rowsAffected != 1 {
			t.Errorf("Expected rows affected to be 1, got %d", rowsAffected)
		}
	})

	t.Run("should return an error when given an id that doesn't exist", func(t *testing.T) {
		repositoryTestsSetup()
		defer repositoryTestsTeardown()
		repository.Db.Exec(queries.ClearTodoTable)
		_, err := repository.Delete(1)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
