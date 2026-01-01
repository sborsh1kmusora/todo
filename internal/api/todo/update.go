package todo

import (
	"errors"
	"log/slog"
	"net/http"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

func (a *api) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseIdFromPath(r)
	if err != nil {
		a.log.Error("error parsing id", slog.Int("id", id), slog.Any("error", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var task model.Task
	if err := parseReqBody(r, &task); err != nil {
		a.log.Error("Error decoding update request", slog.Any("error", err))
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := a.serv.Update(ctx, id, task); err != nil {
		switch {
		case errors.Is(err, appErrors.ErrTaskNotFound):
			a.log.Warn("task not found", slog.Any("task", task))
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, appErrors.ErrTaskTitleIsEmpty):
			a.log.Warn("task title is empty", slog.Any("task", task))
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			a.log.Error("error creating task", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
