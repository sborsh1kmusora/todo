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
	"github.com/sborsh1kmusora/todo/internal/model"
)

func TestAPI_list(t *testing.T) {
	trueVal := true

	tests := []struct {
		name           string
		query          string
		mockSetup      func(m *mocks.MockService)
		expectedStatus int
		expectedBody   []model.Task
	}{
		{
			name:  "default pagination",
			query: "",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					List(gomock.Any(), model.TaskFilter{
						IsDone: nil,
						Limit:  20,
						Offset: 0,
					}).
					Return([]model.Task{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []model.Task{},
		},
		{
			name:  "isDone true",
			query: "?isDone=true&limit=10&offset=5",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					List(gomock.Any(), model.TaskFilter{
						IsDone: &trueVal,
						Limit:  10,
						Offset: 5,
					}).
					Return([]model.Task{
						{ID: 1, Title: "task"},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []model.Task{
				{ID: 1, Title: "task"},
			},
		},
		{
			name:           "invalid isDone",
			query:          "?isDone=abc",
			mockSetup:      func(m *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid limit",
			query:          "?limit=abc",
			mockSetup:      func(m *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid offset",
			query:          "?offset=-1",
			mockSetup:      func(m *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "limit greater than max",
			query: "?limit=1000",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					List(gomock.Any(), model.TaskFilter{
						IsDone: nil,
						Limit:  100,
						Offset: 0,
					}).
					Return([]model.Task{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []model.Task{},
		},
		{
			name:  "service error",
			query: "?limit=10",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error"))
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

			req := httptest.NewRequest(
				http.MethodGet,
				"/todos"+tt.query,
				nil,
			)
			rec := httptest.NewRecorder()

			a.list(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != nil {
				require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

				var got []model.Task
				err := json.NewDecoder(bytes.NewBuffer(rec.Body.Bytes())).Decode(&got)
				require.NoError(t, err)

				require.Equal(t, tt.expectedBody, got)
			}
		})
	}
}
