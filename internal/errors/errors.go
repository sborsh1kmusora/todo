package errors

import "errors"

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrTaskTitleIsEmpty = errors.New("task title is empty")
	ErrTaskAlreadyExist = errors.New("task already exist")
)
