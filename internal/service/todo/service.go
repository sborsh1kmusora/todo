package todo

import (
	"context"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

//go:generate mockgen -source=service.go -destination=mocks/repo_mock.go -package=mocks
type Repository interface {
	Create(context.Context, model.Task) error
	Get(context.Context, int) (model.Task, error)
	List(context.Context, model.TaskFilter) ([]model.Task, error)
	Update(context.Context, int, model.Task) error
	Delete(context.Context, int) error
	Exists(context.Context, int) bool
}

type service struct {
	repo Repository
}

func New(repo Repository) service {
	return service{repo}
}

func (s *service) Create(ctx context.Context, task model.Task) error {
	if s.repo.Exists(ctx, task.ID) {
		return appErrors.ErrTaskAlreadyExist
	}

	if err := validate(task); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return err
	}

	return nil
}

func (s *service) Get(ctx context.Context, id int) (model.Task, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) List(ctx context.Context, filter model.TaskFilter) ([]model.Task, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) Update(ctx context.Context, id int, task model.Task) error {
	if err := validate(task); err != nil {
		return err
	}

	return s.repo.Update(ctx, id, task)
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func validate(task model.Task) error {
	if task.ID <= 0 {
		return appErrors.ErrInvalidTask
	}

	if task.Title == "" {
		return appErrors.ErrTaskTitleIsEmpty
	}

	return nil
}
