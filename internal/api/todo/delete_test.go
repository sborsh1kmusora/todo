package todo

import (
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
)

func TestAPI_delete(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(m *mocks.MockService)
		expectedStatus int
	}{
		{
			name: "success",
			id:   "1",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Delete(gomock.Any(), 1).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "task not found",
			id:   "42",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Delete(gomock.Any(), 42).
					Return(appErrors.ErrTaskNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			id:   "1",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Delete(gomock.Any(), 1).
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
			tt.mockSetup(mockService)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			a := New(mockService, logger)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /todos/{id}", a.delete)

			req := httptest.NewRequest(
				http.MethodDelete,
				"/todos/"+tt.id,
				nil,
			)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
