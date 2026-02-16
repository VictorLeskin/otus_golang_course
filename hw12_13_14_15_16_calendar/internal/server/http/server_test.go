package internalhttp

import (
	"bytes"
	"calendar/internal/app"
	"calendar/internal/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
			}, mockApp)

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
