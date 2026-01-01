package todo

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
)

func (a *api) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseIdFromPath(r)
	if err != nil {
		a.log.Error("error parsing id", slog.Int("id", id), slog.Any("error", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := a.serv.Get(ctx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrTaskNotFound) {
			a.log.Warn("task not found")
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		a.log.Error("failed to get task", slog.Any("id", id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		a.log.Error("error encoding tasks", slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
