package todo

import (
	"bytes"
	"encoding/json"
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

func TestAPI_get(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(m *mocks.MockService)
		expectedStatus int
		expectedBody   *model.Task
	}{
		{
			name:           "invalid id",
			id:             "abc",
			mockSetup:      func(m *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			id:   "1",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Get(gomock.Any(), 1).
					Return(model.Task{ID: 1, Title: "test"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &model.Task{
				ID:    1,
				Title: "test",
			},
		},
		{
			name: "task not found",
			id:   "42",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Get(gomock.Any(), 42).
					Return(model.Task{}, appErrors.ErrTaskNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			id:   "1",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Get(gomock.Any(), 1).
					Return(model.Task{}, errors.New("db error"))
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
			mux.HandleFunc("GET /todos/{id}", a.get)

			req := httptest.NewRequest(
				http.MethodGet,
				"/todos/"+tt.id,
				nil,
			)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != nil {
				require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

				var got model.Task
				err := json.NewDecoder(bytes.NewBuffer(rec.Body.Bytes())).Decode(&got)
				require.NoError(t, err)

				require.Equal(t, *tt.expectedBody, got)
			}
		})
	}
}
