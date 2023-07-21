package todo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/raphael-foliveira/fiber-todo/pkg/database"
)

var db *database.Database
var app *fiber.App

func TestMain(m *testing.M) {
	fmt.Println("running tests...")
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db = database.MustGetDatabase(os.Getenv("TEST_DATABASE_URL"))
	db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	db.Migrate("../database/schema.sql")
	db.Exec("INSERT INTO todo (title, description, completed) VALUES ('test', 'test', false), ('test2', 'test2', false)")
	app = fiber.New()
	group := app.Group("/todos")
	GetTodoRoutes(group, db.DB)

	code := m.Run()

	db.Exec("DELETE FROM todo")

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	var todo TodoDto
	faker.FakeData(&todo)
	todoW := new(bytes.Buffer)
	err := json.NewEncoder(todoW).Encode(&todo)
	if err != nil {
		t.Errorf("Error marshalling todo: %v", err)
	}
	req, err := http.NewRequest("POST", "/todos", todoW)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v", http.StatusCreated, res.StatusCode)
	}
}

func TestList(t *testing.T) {
	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		fmt.Println(err)
	}
	res, err := app.Test(req)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Expected status code 200, got", res.StatusCode)
	}
}

func TestRetrieve(t *testing.T) {
	req, err := http.NewRequest("GET", "/todos/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	res, err := app.Test(req)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Expected status code 200, got", res.StatusCode)
	}
}

func TestUpdate(t *testing.T) {
	var todo Todo
	faker.FakeData(&todo)
	todoW := new(bytes.Buffer)
	err := json.NewEncoder(todoW).Encode(&todo)
	if err != nil {
		t.Errorf("Error marshalling todo: %v", err)
	}
	req, err := http.NewRequest("PUT", "/todos/1", todoW)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, res.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/todos/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	res, err := app.Test(req)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode != http.StatusNoContent {
		fmt.Println("Expected status code 200, got", res.StatusCode)
	}
}
