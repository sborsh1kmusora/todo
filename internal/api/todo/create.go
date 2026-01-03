package todo

import (
	"errors"
	"log/slog"
	"net/http"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

// CreateTask godoc
// @Summary Создать задачу
// @Description Создаёт новую задачу
// @Tags todos
// @Accept json
// @Produce json
// @Param task body model.Task true "Задача для создания"
// @Success 201 "Задача успешно создана"
// @Failure 400 {object} errors.ErrorResponse "Невалидное тело запроса или пустой title"
// @Failure 409 {object} errors.ErrorResponse "Задача уже существует"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /todos [post]
func (a *api) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var task model.Task
	if err := a.parseReqBody(r, &task); err != nil {
		a.log.Error("error decoding task", slog.Any("error", err))
		appErrors.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.serv.Create(ctx, task); err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidTask):
			a.log.Warn("received invalid task", slog.Any("task", task))
			appErrors.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, appErrors.ErrTaskAlreadyExist):
			a.log.Warn("task already exists", slog.Any("task", task))
			appErrors.WriteError(w, http.StatusConflict, err.Error())
		case errors.Is(err, appErrors.ErrTaskTitleIsEmpty):
			a.log.Warn("task title is empty", slog.Any("task", task))
			appErrors.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			a.log.Error("error creating task", slog.Any("error", err))
			appErrors.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	a.log.Info("task created successfully", slog.Any("task", task))

	w.WriteHeader(http.StatusCreated)
}
