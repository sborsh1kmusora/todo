package todo

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sborsh1kmusora/todo/internal/model"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

func (a *api) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query()
	isDoneStr := q.Get("isDone")

	var isDone *bool
	if isDoneStr != "" {
		v, err := strconv.ParseBool(isDoneStr)
		if err != nil {
			a.log.Error("failed to parse isDone query parameter", slog.String("isDone", isDoneStr))
			http.Error(w, "invalid isDone value", http.StatusBadRequest)
			return
		}
		isDone = &v
	}

	limit, offset, err := parsePagination(q)
	if err != nil {
		a.log.Error("failed to parse pagination", slog.String("error", err.Error()))
		http.Error(w, "invalid pagination values", http.StatusBadRequest)
		return
	}

	filter := model.TaskFilter{
		IsDone: isDone,
		Limit:  limit,
		Offset: offset,
	}

	tasks, err := a.serv.List(ctx, filter)
	if err != nil {
		a.log.Error("error listing tasks", slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	a.log.Info("tasks listed", slog.Any("tasks", tasks))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		a.log.Error("error encoding tasks", slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func parsePagination(q url.Values) (limit, offset int, err error) {
	limit = defaultLimit
	offset = 0

	if v := q.Get("limit"); v != "" {
		limit, err = strconv.Atoi(v)
		if err != nil || limit <= 0 {
			return 0, 0, errors.New("invalid limit")
		}
		if limit > maxLimit {
			limit = maxLimit
		}
	}

	if v := q.Get("offset"); v != "" {
		offset, err = strconv.Atoi(v)
		if err != nil || offset < 0 {
			return 0, 0, errors.New("invalid offset")
		}
	}

	return limit, offset, nil
}
