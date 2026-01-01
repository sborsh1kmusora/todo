package errors

import "errors"

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrInvalidTask      = errors.New("invalid task")
	ErrTaskTitleIsEmpty = errors.New("task title is empty")
	ErrTaskAlreadyExist = errors.New("task already exists")
)
