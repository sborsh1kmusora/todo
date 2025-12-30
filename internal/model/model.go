package model

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	IsDone bool   `json:"is_done"`
}
