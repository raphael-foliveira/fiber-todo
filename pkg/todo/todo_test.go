package todo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
)

// var db *database.Database
var app *fiber.App
var mr *mockRepository

type todoTest struct {
	name         string
	modifier     func(*bytes.Buffer)
	urlFunc      func() string
	expectStatus int
}

type mockRepository struct {
	todos      []Todo
	shouldFail bool
}

func (mr *mockRepository) Create(todo TodoDto) (Todo, error) {
	id := 0
	for _, t := range mr.todos {
		if t.Title == todo.Title {
			return Todo{}, errors.New("todo already exists in mock repository")
		}
		id++
	}
	mr.todos = append(mr.todos, Todo{
		Id:          id,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	})
	return Todo{}, nil
}

func (mr *mockRepository) List() ([]Todo, error) {
	if mr.shouldFail {
		return nil, errors.New("error listing todos")
	}
	return mr.todos, nil
}

func (mr *mockRepository) Retrieve(id int) (Todo, error) {
	for _, todo := range mr.todos {
		if todo.Id == id {
			return todo, nil
		}
	}
	return Todo{}, errors.New("todo not found in mock repository")
}

func (mr *mockRepository) Update(todo Todo) (Todo, error) {
	if mr.shouldFail {
		return Todo{}, errors.New("error updating todo")
	}
	for _, t := range mr.todos {
		if t.Id == todo.Id {
			fmt.Println(todo)
			t.Title = todo.Title
			t.Description = todo.Description
			t.Completed = todo.Completed
			return t, nil
		}
	}
	return Todo{Id: 0}, nil
}

func (mr *mockRepository) Delete(id int) (int64, error) {
	if mr.shouldFail {
		return 0, errors.New("error deleting todo")
	}
	for i, t := range mr.todos {
		if t.Id == id {
			mr.todos = append(mr.todos[:i], mr.todos[i+1:]...)
			return 1, nil
		}
	}
	return 0, nil
}

func (mr *mockRepository) InsertFixtures() {
	mr.todos = []Todo{}
	for i := 0; i < 30; i++ {
		var todo Todo
		faker.FakeData(&todo)
		todo.Id = i + 1
		mr.todos = append(mr.todos, todo)
	}
}

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
func todoTestsSetup() {
	// db = database.MustGetDatabase(config.Database.Url)
	// db.Exec(queries.RecreateSchema)
	// db.Migrate()
	// db.Exec(queries.InsertTodoFixtures)
	app = fiber.New()
	group := app.Group("/todos")
	mr = new(mockRepository)
	controller := NewTodoController(mr)
	mr.InsertFixtures()
	GetTodoRoutes(group, controller)
}

func todoTestsTeardown() {
	mr.todos = []Todo{}
	mr.shouldFail = false
}

func TestMain(m *testing.M) {
	fmt.Println("running tests...")
	os.Exit(m.Run())
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
				existingTodoTitle := mr.todos[0].Title
				todo := Todo{
					Title:       existingTodoTitle,
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
			todoTestsSetup()
			defer todoTestsTeardown()
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
				mr.shouldFail = true
			},
			func() string { return "/todos" },
			500,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			todoTestsSetup()
			defer todoTestsTeardown()
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
			todoTestsSetup()
			defer todoTestsTeardown()
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
			func(b *bytes.Buffer) {
			},
			func() string {
				nonExistingTodo := mr.todos[0]
				mr.todos = mr.todos[1:]
				url := fmt.Sprintf("/todos/%v", nonExistingTodo.Id)
				fmt.Println(url)
				return url
			},
			404,
		},
		{
			"test update fail",
			func(b *bytes.Buffer) {
				mr.shouldFail = true
			},
			func() string { return "/todos/1" },
			500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			todoTestsSetup()
			defer todoTestsTeardown()
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
			"test delete",
			func(b *bytes.Buffer) {},
			func() string { return "/todos/1" },
			204,
		},
		{
			"test delete invalid",
			func(b *bytes.Buffer) {},
			func() string { return "/todos/invalid" },
			422,
		},
		{
			"test delete non existing",
			func(b *bytes.Buffer) {},
			func() string {
				nonExistingTodo := mr.todos[0]
				mr.todos = mr.todos[1:]
				url := fmt.Sprintf("/todos/%v", nonExistingTodo.Id)
				return url
			},
			404,
		},
		{
			"test delete fail",
			func(b *bytes.Buffer) {
				mr.shouldFail = true
			},
			func() string { return "/todos/1" },
			500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			todoTestsSetup()
			defer todoTestsTeardown()
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
