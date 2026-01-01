package todo

import (
	"errors"
	"log/slog"
	"net/http"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
)

func (a *api) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseIdFromPath(r)
	if err != nil {
		a.log.Error("error parsing id", slog.Int("id", id), slog.Any("error", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := a.serv.Delete(ctx, id); err != nil {
		if errors.Is(err, appErrors.ErrTaskNotFound) {
			a.log.Warn("task not found")
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		a.log.Error("failed to delete task", slog.Any("id", id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	a.log.Info("successfully deleted task", slog.Any("id", id))

	w.WriteHeader(http.StatusNoContent)
}
