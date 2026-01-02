package todo

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sborsh1kmusora/todo/internal/model"
)

//go:generate mockgen -source=api.go -destination=mocks/service_mock.go -package=mocks
type Service interface {
	Create(context.Context, model.Task) error
	Get(context.Context, int) (model.Task, error)
	List(context.Context, model.TaskFilter) ([]model.Task, error)
	Update(context.Context, int, model.Task) error
	Delete(context.Context, int) error
}

type api struct {
	serv Service
	log  *slog.Logger
}

func New(serv Service, log *slog.Logger) api {
	return api{
		serv: serv,
		log:  log,
	}
}

func (a *api) Todos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.create(w, r)
	case http.MethodGet:
		a.list(w, r)
	default:
		a.log.Error("Method not allowed")
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *api) TodoById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.get(w, r)
	case http.MethodPut:
		a.update(w, r)
	case http.MethodDelete:
		a.delete(w, r)
	default:
		a.log.Error("method not allowed")
		w.Header().Set("Allow", "GET, PUT, DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *api) parseReqBody(r *http.Request, dst *model.Task) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			a.log.Error("failed to close body")
		}
	}(r.Body)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(dst)
}

func parseIdFromPath(r *http.Request) (int, error) {
	idStr := r.PathValue("id")

	return strconv.Atoi(idStr)
}
