package internalhttp

import (
	"bytes"
	"calendar/internal/app"
	"calendar/internal/logger"
	memorystorage "calendar/internal/storage/memory"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_Response(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		response string
	}{
		{
			"hello",
			"/hello",
			"Hello from Calendar Service!\n",
		},
		{
			"root",
			"/",
			"Hello, World!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logBuffer := &bytes.Buffer{}
			log := logger.NewWriterLogger("debug", logBuffer)

			mockApp := &app.App{
				Logger: log,
			}

			srv := NewServer(Config{
				Host: "localhost",
				Port: "8080",
			}, mockApp, nil)

			srv.RegisterHandlers()

			req := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			// Получаем handler и вызываем его
			handler := srv.GetHandler() // нужно добавить этот метод
			handler.ServeHTTP(w, req)

			// Проверяем ответ
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, string(body), tt.response)
			assert.True(t, strings.HasPrefix(logBuffer.String(), "[D] Request started"))
			assert.Contains(t, logBuffer.String(), "[I] Request completed")
		})
	}
}

func TestCreateEvent_Success(t *testing.T) {
	// Сетап с реальным in-memory storage
	store := memorystorage.New()

	logBuffer := &bytes.Buffer{}
	log := logger.NewWriterLogger("debug", logBuffer)

	mockApp := &app.App{
		Logger: log,
	}

	server := NewServer(Config{
		Host: "localhost",
		Port: "8080",
	}, mockApp, store)

	// Регистрируем хендлер
	mux := http.NewServeMux()
	mux.HandleFunc("POST /events", server.handleCreateEvent)
	handler := LoggingMiddleware(log)(mux)

	// Подготовка запроса
	reqBody := CreateEventRequest{
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   time.Date(2026, 2, 25, 15, 30, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 2, 25, 16, 30, 0, 0, time.UTC),
		UserID:      "user-123",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/events", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp EventResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "Test Event", resp.Title)

	// Проверяем, что в storage действительно появилась запись
	savedEvent, err := store.GetEvent(context.Background(), resp.ID)
	assert.NoError(t, err)
	assert.Equal(t, resp.Title, savedEvent.Title)
}
