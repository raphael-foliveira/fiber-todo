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

type todoTest struct {
	name         string
	modifier     func(*bytes.Buffer)
	urlFunc      func() string
	expectStatus int
}

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
	tests := []todoTest{
		{
			name:         "create valid todo",
			expectStatus: 201,
			urlFunc:      func() string { return "/todos" },
			modifier:     func(todoW *bytes.Buffer) {},
		},
		{
			name:         "create invalid todo",
			expectStatus: 400,
			urlFunc:      func() string { return "/todos" },
			modifier: func(b *bytes.Buffer) {
				b.Reset()
				b.WriteString("invalid")
			},
		},
		{
			name:         "create conflicting todo",
			expectStatus: 409,
			urlFunc:      func() string { return "/todos" },
			modifier: func(b *bytes.Buffer) {
				b.Reset()
				todo := Todo{
					Title:       "test2",
					Description: "Test",
					Completed:   false,
				}
				err := json.NewEncoder(b).Encode(&todo)
				if err != nil {
					t.Errorf("Error encoding todo: %v", err)
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setup()
			defer teardown()
			todoW, err := createTodoBodyHelper()
			if err != nil {
				t.Errorf("Error creating todo: %v", err)
			}
			test.modifier(todoW)
			req, err := http.NewRequest("POST", test.urlFunc(), todoW)
			if err != nil {
				t.Errorf("Error creating request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			if err != nil {
				t.Errorf("Error sending request: %v", err)
			}
			if res.StatusCode != test.expectStatus {
				t.Errorf("Expected status code %v, got %v", test.expectStatus, res.StatusCode)
			}
		})
	}
}

func TestList(t *testing.T) {
	tests := []todoTest{
		{
			name:         "test list",
			expectStatus: 200,
			modifier:     func(b *bytes.Buffer) {},
			urlFunc:      func() string { return "/todos" },
		},
		{
			name:         "test list fail",
			expectStatus: 500,
			modifier: func(b *bytes.Buffer) {
				db.Close()
			},
			urlFunc: func() string { return "/todos" },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setup()
			defer teardown()
			test.modifier(new(bytes.Buffer))
			req, err := http.NewRequest("GET", test.urlFunc(), nil)
			if err != nil {
				t.Error(err)
			}
			res, err := app.Test(req)
			if err != nil {
				t.Error(err)
			}
			if res.StatusCode != test.expectStatus {
				t.Errorf("Expected status code %v, got %v", test.expectStatus, res.StatusCode)
			}
		})
	}
}

func TestRetrieve(t *testing.T) {
	var tests = []todoTest{
		{
			name:         "test retrieve",
			modifier:     func(b *bytes.Buffer) {},
			expectStatus: 200,
			urlFunc:      func() string { return "/todos/1" },
		},
		{
			name:         "test retrieve invalid",
			modifier:     func(b *bytes.Buffer) {},
			expectStatus: 422,
			urlFunc:      func() string { return "/todos/invalid" },
		},
		{
			name:         "test retrieve non existing",
			modifier:     func(b *bytes.Buffer) {},
			expectStatus: 404,
			urlFunc:      func() string { return "/todos/999" },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setup()
			defer teardown()
			b := new(bytes.Buffer)
			test.modifier(b)
			req, err := http.NewRequest("GET", test.urlFunc(), nil)
			if err != nil {
				t.Error(err)
			}
			res, err := app.Test(req)
			if err != nil {
				t.Error(err)
			}
			if res.StatusCode != test.expectStatus {
				t.Errorf("Expected status code %v, got %v", test.expectStatus, res.StatusCode)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []todoTest{
		{
			name: "test update",
			modifier: func(b *bytes.Buffer) {
			},
			expectStatus: 200,
			urlFunc:      func() string { return "/todos/1" },
		},
		{
			name: "test update invalid",
			modifier: func(b *bytes.Buffer) {
				b.Reset()
				b.WriteString("invalid")
			},
			expectStatus: 400,
			urlFunc:      func() string { return "/todos/1" },
		},
		{
			name:         "test update invalid id",
			modifier:     func(b *bytes.Buffer) {},
			expectStatus: 422,
			urlFunc:      func() string { return "/todos/invalid" },
		},
		{
			name:         "test update non existing",
			modifier:     func(b *bytes.Buffer) {},
			expectStatus: 404,
			urlFunc:      func() string { return "/todos/999" },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setup()
			defer teardown()
			todoW, err := createTodoBodyHelper()
			if err != nil {
				t.Errorf("Error creating todo: %v", err)
			}
			test.modifier(todoW)
			req, err := http.NewRequest("PUT", test.urlFunc(), todoW)
			if err != nil {
				t.Errorf("Error creating request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)
			if err != nil {
				t.Errorf("Error sending request: %v", err)
			}

			if res.StatusCode != test.expectStatus {
				t.Errorf("Expected status code %v, got %v", test.expectStatus, res.StatusCode)
			}
		})
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
	tests := []todoTest{
		{
			name:         "test delete",
			modifier:     func(b *bytes.Buffer) {},
			urlFunc:      func() string { return "/todos/1" },
			expectStatus: 204,
		},
		{
			name:         "test delete invalid",
			modifier:     func(b *bytes.Buffer) {},
			urlFunc:      func() string { return "/todos/invalid" },
			expectStatus: 422,
		},
		{
			name:         "test delete non existing",
			modifier:     func(b *bytes.Buffer) {},
			urlFunc:      func() string { return "/todos/999" },
			expectStatus: 404,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setup()
			defer teardown()
			b := new(bytes.Buffer)
			test.modifier(b)
			req, err := http.NewRequest("DELETE", test.urlFunc(), nil)
			if err != nil {
				t.Error(err)
			}
			res, err := app.Test(req)
			if err != nil {
				t.Error(err)
			}
			if res.StatusCode != test.expectStatus {
				t.Errorf("Expected status code %v, got %v", test.expectStatus, res.StatusCode)
			}
		})
	}
}
