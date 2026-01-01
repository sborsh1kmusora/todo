package todo

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (a *api) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := a.serv.List(ctx)
	if err != nil {
		a.log.Error("error listing tasks", slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	a.log.Info("tasks listed", slog.Any("tasks", tasks))

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		a.log.Error("error encoding tasks", slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
