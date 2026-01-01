package todo

import (
	"errors"
	"log/slog"
	"net/http"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

func (a *api) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var task model.Task
	if err := parseReqBody(r, &task); err != nil {
		a.log.Error("Error decoding task", slog.Any("error", err))
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := a.serv.Create(ctx, task); err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidTask):
			a.log.Warn("received invalid task", slog.Any("task", task))
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, appErrors.ErrTaskAlreadyExist):
			a.log.Warn("task already exists", slog.Any("task", task))
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, appErrors.ErrTaskTitleIsEmpty):
			a.log.Warn("task title is empty", slog.Any("task", task))
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			a.log.Error("error creating task", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	a.log.Info("task created successfully", slog.Any("task", task))

	w.WriteHeader(http.StatusCreated)
}
