package todo

import (
	"bytes"
	"errors"
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

func TestAPI_create(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockSetup      func(m *mocks.MockService)
		expectedStatus int
	}{
		{
			name:           "invalid json",
			body:           `{"id":`,
			mockSetup:      func(m *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			body: `{"id":1,"title":"test"}`,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Create(gomock.Any(), model.Task{ID: 1, Title: "test"}).
					Return(nil).
					Times(1)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid task",
			body: `{"id":0,"title":"test"}`,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(appErrors.ErrInvalidTask)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty title",
			body: `{"id":1,"title":""}`,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(appErrors.ErrTaskTitleIsEmpty)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "task already exists",
			body: `{"id":1,"title":"test"}`,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(appErrors.ErrTaskAlreadyExist)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "internal error",
			body: `{"id":1,"title":"test"}`,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("unexpected error"))
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

			logger := slog.New(
				slog.NewTextHandler(bytes.NewBuffer(nil), nil),
			)

			a := New(mockService, logger)

			req := httptest.NewRequest(
				http.MethodPost,
				"/todos",
				bytes.NewBufferString(tt.body),
			)
			rec := httptest.NewRecorder()

			a.create(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
