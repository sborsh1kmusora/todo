package model

import "time"

type Task struct {
	ID        int       `json:"id" example:"1"`
	Title     string    `json:"title" example:"Написать web сервис"`
	Desc      string    `json:"desc" example:"Реализовать CRUD операции"`
	IsDone    bool      `json:"is_done" example:"false"`
	CreatedAt time.Time `json:"created_at" example:"2026-01-01"`
}

type TaskFilter struct {
	IsDone *bool
	Limit  int
	Offset int
}
