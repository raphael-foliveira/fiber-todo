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
	"github.com/raphael-foliveira/fiber-todo/pkg/common"
	"github.com/raphael-foliveira/fiber-todo/pkg/database"
	"github.com/raphael-foliveira/fiber-todo/pkg/database/queries"
)

var db *database.Database
var app *fiber.App
var config common.Config

func createTodoBodyHelper() (*bytes.Buffer, error) {
	var todo Todo
	faker.FakeData(&todo)
	todoW := new(bytes.Buffer)
	err := json.NewEncoder(todoW).Encode(&todo)
	if err != nil {
		return nil, err
	}
	return todoW, err
}
func setup() {
	db = database.MustGetDatabase(config.Database.Url)
	db.Exec(queries.RecreateSchema)
	db.Migrate()
	db.Exec(queries.InsertTodoFixtures)
	app = fiber.New()
	group := app.Group("/todos")
	GetTodoRoutes(group, db.DB)
}

func teardown() {
	db.Exec(queries.ClearTodoTable)
}

func TestMain(m *testing.M) {
	config = common.ReadTestCfg()
	fmt.Println("running tests...")
	code := m.Run()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	setup()
	defer teardown()
	todoW, err := createTodoBodyHelper()
	if err != nil {
		t.Errorf("Error creating todo: %v", err)
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

func TestCreateInvalid(t *testing.T) {
	setup()
	defer teardown()
	todoW, err := createTodoBodyHelper()
	if err != nil {
		t.Errorf("Error creating todo: %v", err)
	}
	todoW.WriteString("invalid")
	req, err := http.NewRequest("POST", "/todos", todoW)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestCreateConflict(t *testing.T) {
	setup()
	defer teardown()
	todo := Todo{
		Title:       "test2",
		Description: "test2",
		Completed:   false,
	}
	todoW, err := json.Marshal(todo)
	if err != nil {
		t.Errorf("Error creating todo: %v", err)
	}
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(todoW))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code %v, got %v", http.StatusConflict, res.StatusCode)
	}
}

func TestList(t *testing.T) {
	setup()
	defer teardown()
	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
}

func TestListFail(t *testing.T) {
	setup()
	defer teardown()
	db.Close()
	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", res.StatusCode)
	}
}

func TestRetrieve(t *testing.T) {
	setup()
	defer teardown()
	req, err := http.NewRequest("GET", "/todos/1", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
}

func TestRetrieveInvalidId(t *testing.T) {
	setup()
	defer teardown()
	req, err := http.NewRequest("GET", "/todos/invalid", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422, got %d", res.StatusCode)
	}
}

func TestRetrieveNotFound(t *testing.T) {
	setup()
	defer teardown()
	req, err := http.NewRequest("GET", "/todos/999", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", res.StatusCode)
	}
}

func TestUpdate(t *testing.T) {
	setup()
	defer teardown()
	todoW, err := createTodoBodyHelper()
	if err != nil {
		t.Errorf("Error creating todo: %v", err)
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

func TestUpdateInvalid(t *testing.T) {
	setup()
	defer teardown()
	todoW, err := createTodoBodyHelper()
	if err != nil {
		t.Errorf("Error creating todo: %v", err)
	}
	todoW.WriteString("invalid")
	req, err := http.NewRequest("PUT", "/todos/1", todoW)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestUpdateInvalidId(t *testing.T) {
	setup()
	defer teardown()
	todoW, err := createTodoBodyHelper()
	if err != nil {
		t.Errorf("Error creating todo: %v", err)
	}
	req, err := http.NewRequest("PUT", "/todos/invalid", todoW)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %v, got %v", http.StatusUnprocessableEntity, res.StatusCode)
	}
}

func TestUpdateNonExisting(t *testing.T) {
	setup()
	defer teardown()
	todoW, err := createTodoBodyHelper()
	if err != nil {
		t.Errorf("Error creating todo: %v", err)
	}
	req, err := http.NewRequest("PUT", "/todos/999", todoW)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %v, got %v", http.StatusNotFound, res.StatusCode)
	}
}

func TestUpdateFail(t *testing.T) {
	setup()
	defer teardown()
	db.Close()
	todoW, err := createTodoBodyHelper()
	if err != nil {
		t.Errorf("Error creating todo: %v", err)
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

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %v, got %v", http.StatusInternalServerError, res.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	setup()
	defer teardown()
	req, err := http.NewRequest("DELETE", "/todos/1", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
}

func TestDeleteInvalidId(t *testing.T) {
	setup()
	defer teardown()
	req, err := http.NewRequest("DELETE", "/todos/invalid", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422, got %d", res.StatusCode)
	}
}

func TestDeleteNotFound(t *testing.T) {
	setup()
	defer teardown()
	req, err := http.NewRequest("DELETE", "/todos/999", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", res.StatusCode)
	}
}
