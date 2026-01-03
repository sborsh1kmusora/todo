package todo

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// ListTasks godoc
// @Summary Получить список задач
// @Description Возвращает список задач с фильтрацией и пагинацией
// @Tags todos
// @Accept json
// @Produce json
// @Param isDone query bool false "Фильтр по статусу выполнения задачи"
// @Param limit query int false "Количество задач (по умолчанию 20, максимум 100)" minimum(1) maximum(100)
// @Param offset query int false "Смещение (offset) для пагинации" minimum(0)
// @Success 200 {array} model.Task
// @Failure 400 {object} errors.ErrorResponse "Некорректные query-параметры"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /todos [get]
func (a *api) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query()
	isDoneStr := q.Get("isDone")

	var isDone *bool
	if isDoneStr != "" {
		v, err := strconv.ParseBool(isDoneStr)
		if err != nil {
			a.log.Error("failed to parse isDone query parameter", slog.String("isDone", isDoneStr))
			appErrors.WriteError(w, http.StatusBadRequest, "invalid isDone value")
			return
		}
		isDone = &v
	}

	limit, offset, err := parsePagination(q)
	if err != nil {
		a.log.Error("failed to parse pagination", slog.String("error", err.Error()))
		appErrors.WriteError(w, http.StatusBadRequest, "invalid pagination values")
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
		appErrors.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	a.log.Info("tasks listed", slog.Any("tasks", tasks))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		a.log.Error("error encoding tasks", slog.Any("error", err))
		appErrors.WriteError(w, http.StatusInternalServerError, "internal server error")
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
