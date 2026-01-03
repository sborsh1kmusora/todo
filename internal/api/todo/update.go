package todo

import (
	"errors"
	"log/slog"
	"net/http"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

// UpdateTask godoc
// @Summary Обновить задачу
// @Description Обновляет существующую задачу по ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param task body model.Task true "Данные для обновления задачи"
// @Success 204 "Задача успешно обновлена"
// @Failure 400 {object} errors.ErrorResponse "Некорректный ID или тело запроса"
// @Failure 404 {object} errors.ErrorResponse "Задача не найдена"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /todos/{id} [put]
func (a *api) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseIdFromPath(r)
	if err != nil {
		a.log.Error("error parsing id", slog.Int("id", id), slog.Any("error", err))
		appErrors.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var task model.Task
	if err := a.parseReqBody(r, &task); err != nil {
		a.log.Error("error decoding update request", slog.Any("error", err))
		appErrors.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := a.serv.Update(ctx, id, task); err != nil {
		switch {
		case errors.Is(err, appErrors.ErrTaskNotFound):
			a.log.Warn("task not found", slog.Any("task", task))
			http.Error(w, err.Error(), http.StatusNotFound)
			appErrors.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, appErrors.ErrTaskTitleIsEmpty):
			a.log.Warn("task title is empty", slog.Any("task", task))
			appErrors.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			a.log.Error("error creating task", slog.Any("error", err))
			appErrors.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
