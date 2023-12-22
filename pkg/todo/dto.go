package todo

type CreateTodoDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type UpdateTodoDto CreateTodoDto

type CreateResponse struct {
	Id int `json:"id"`
}
