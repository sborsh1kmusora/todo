package todo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

type Service interface {
	Create(context.Context, model.Task) error
	Get(context.Context, int) (model.Task, error)
	List(context.Context) ([]model.Task, error)
	Update(context.Context, int, model.Task) error
	Delete(context.Context, int) error
}

type api struct {
	serv Service
}

func New(s Service) api {
	return api{
		serv: s,
	}
}

func (a *api) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := a.serv.Create(ctx, task); err != nil {
		switch {
		case errors.Is(err, appErrors.ErrTaskTitleIsEmpty):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, appErrors.ErrTaskAlreadyExist):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
