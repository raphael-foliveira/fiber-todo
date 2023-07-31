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
			"create valid todo",
			func(b *bytes.Buffer) {},
			func() string { return "/todos" },
			201,
		},
		{
			"create invalid todo",
			func(b *bytes.Buffer) {
				b.Reset()
				b.WriteString("invalid")
			},
			func() string { return "/todos" },
			400,
		},
		{
			"create conflicting todo",
			func(b *bytes.Buffer) {
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
			func() string { return "/todos" },
			409,
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
			"test list",
			func(b *bytes.Buffer) {},
			func() string { return "/todos" },
			200,
		},
		{
			"test list fail",
			func(b *bytes.Buffer) {
				db.Close()
			},
			func() string { return "/todos" },
			500,
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
			"test retrieve",
			func(b *bytes.Buffer) {},
			func() string { return "/todos/1" },
			200,
		},
		{
			"test retrieve invalid",
			func(b *bytes.Buffer) {},
			func() string { return "/todos/invalid" },
			422,
		},
		{
			"test retrieve non existing",
			func(b *bytes.Buffer) {},
			func() string { return "/todos/999" },
			404,
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
			"test update",
			func(b *bytes.Buffer) {
			},
			func() string { return "/todos/1" },
			200,
		},
		{
			"test update invalid",
			func(b *bytes.Buffer) {
				b.Reset()
				b.WriteString("invalid")
			},
			func() string { return "/todos/1" },
			400,
		},
		{
			"test update invalid id",
			func(b *bytes.Buffer) {},
			func() string { return "/todos/invalid" },
			422,
		},
		{
			"test update non existing",
			func(b *bytes.Buffer) {},
			func() string { return "/todos/999" },
			404,
		},
		{
			"test update fail",
			func(b *bytes.Buffer) {
				db.Close()
			},
			func() string { return "/todos/1" },
			500,
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
		{
			name: "test delete fail",
			modifier: func(b *bytes.Buffer) {
				db.Close()
			},
			urlFunc:      func() string { return "/todos/1" },
			expectStatus: 500,
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
