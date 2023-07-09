package todo

type TodoDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type CreateResponse struct {
	Id int `json:"id"`
}
