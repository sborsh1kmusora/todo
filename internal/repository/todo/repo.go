package todo

import (
	"context"
	"sync"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

type repo struct {
	mu      sync.RWMutex
	storage map[int]model.Task
}

func New() repo {
	return repo{
		storage: make(map[int]model.Task),
		mu:      sync.RWMutex{},
	}
}

func (r *repo) Create(_ context.Context, task model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[task.ID] = task

	return nil
}

func (r *repo) Get(_ context.Context, id int) (model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.storage[id]
	if !ok {
		return model.Task{}, appErrors.ErrTaskNotFound
	}

	return task, nil
}

func (r *repo) List(_ context.Context) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]model.Task, 0, len(r.storage))
	for _, task := range r.storage {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *repo) Update(_ context.Context, id int, task model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[id]; !ok {
		return appErrors.ErrTaskNotFound
	}

	r.storage[id] = task

	return nil
}

func (r *repo) Delete(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.storage, id)

	return nil
}

func (r *repo) Exists(_ context.Context, id int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.storage[id]; ok {
		return true
	}

	return false
}
