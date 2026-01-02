package todo

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sborsh1kmusora/todo/internal/api/todo/mocks"
	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
)

func TestAPI_update(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		body           string
		mockSetup      func(m *mocks.MockService)
		expectedStatus int
	}{
		{
			name:           "invalid id",
			id:             "abc",
			body:           `{"title":"test"}`,
			mockSetup:      func(m *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid body",
			id:             "1",
			body:           `{"title":`,
			mockSetup:      func(m *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			id:   "1",
			body: `{"title":"updated"}`,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Update(
						gomock.Any(),
						1,
						model.Task{Title: "updated"},
					).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "task not found",
			id:   "42",
			body: `{"title":"updated"}`,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Update(gomock.Any(), 42, gomock.Any()).
					Return(appErrors.ErrTaskNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "empty title",
			id:   "1",
			body: `{"title":""}`,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Update(gomock.Any(), 1, gomock.Any()).
					Return(appErrors.ErrTaskTitleIsEmpty)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			id:   "1",
			body: `{"title":"updated"}`,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Update(gomock.Any(), 1, gomock.Any()).
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			a := New(mockService, logger)

			mux := http.NewServeMux()
			mux.HandleFunc("PUT /todos/{id}", a.update)

			req := httptest.NewRequest(
				http.MethodPut,
				"/todos/"+tt.id,
				bytes.NewBufferString(tt.body),
			)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
