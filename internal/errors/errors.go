package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrInvalidTask      = errors.New("invalid task")
	ErrTaskTitleIsEmpty = errors.New("task title is empty")
	ErrTaskAlreadyExist = errors.New("task already exists")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}
