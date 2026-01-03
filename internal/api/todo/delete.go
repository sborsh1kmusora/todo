package todo

import (
	"errors"
	"log/slog"
	"net/http"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
)

// DeleteTask godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по идентификатору
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 204 "Задача успешно удалена"
// @Failure 400 {object} errors.ErrorResponse "Некорректный ID"
// @Failure 404 {object} errors.ErrorResponse "Задача не найдена"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /todos/{id} [delete]
func (a *api) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseIdFromPath(r)
	if err != nil {
		a.log.Error("error parsing id", slog.Int("id", id), slog.Any("error", err))
		appErrors.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := a.serv.Delete(ctx, id); err != nil {
		if errors.Is(err, appErrors.ErrTaskNotFound) {
			a.log.Warn("task not found")
			appErrors.WriteError(w, http.StatusNotFound, "task not found")
			return
		}
		a.log.Error("failed to delete task", slog.Any("id", id), slog.Any("error", err))
		appErrors.WriteError(w, http.StatusNotFound, "internal server error")
		return
	}

	a.log.Info("successfully deleted task", slog.Any("id", id))

	w.WriteHeader(http.StatusNoContent)
}
