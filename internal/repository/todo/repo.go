package todo

import (
	"context"
	"sort"
	"sync"
	"time"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

type repo struct {
	mu      sync.RWMutex
	storage map[int]model.Task
}

func New() repo {
	return repo{
		mu:      sync.RWMutex{},
		storage: make(map[int]model.Task),
	}
}

func (r *repo) Create(_ context.Context, task model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.CreatedAt = time.Now()
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

func (r *repo) List(_ context.Context, filter model.TaskFilter) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]model.Task, 0, len(r.storage))
	for _, task := range r.storage {
		if filter.IsDone != nil && task.IsDone != *filter.IsDone {
			continue
		}
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	if filter.Offset >= len(tasks) {
		return []model.Task{}, nil
	}

	end := filter.Offset + filter.Limit
	if end > len(tasks) {
		end = len(tasks)
	}

	result := make([]model.Task, end-filter.Offset)
	copy(result, tasks[filter.Offset:end])

	return result, nil
}

func (r *repo) Update(_ context.Context, id int, task model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[id]; !ok {
		return appErrors.ErrTaskNotFound
	}

	task.ID = id
	task.CreatedAt = time.Now()

	r.storage[id] = task

	return nil
}

func (r *repo) Delete(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[id]; !ok {
		return appErrors.ErrTaskNotFound
	}

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
