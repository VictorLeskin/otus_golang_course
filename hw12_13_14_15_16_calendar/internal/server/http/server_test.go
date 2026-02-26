package internalhttp

import (
	"bytes"
	"calendar/internal/app"
	"calendar/internal/logger"
	"calendar/internal/storage"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

// MockLogger обёртка над logger.NewWriterLogger для тестов.
type MockLogger struct {
	buf    strings.Builder
	logger *logger.Logger
}

func NewMockLogger() *MockLogger {
	ml := &MockLogger{}
	ml.logger = logger.NewWriterLogger("info", &ml.buf)
	return ml
}

// internal/server/grpc/mock_storage_test.go ...
type MockStorage struct {
	CreateEventFunc func(ctx context.Context, event *storage.Event) error
	UpdateEventFunc func(ctx context.Context, event *storage.Event) error
	DeleteEventFunc func(ctx context.Context, id string) error
	GetEventFunc    func(ctx context.Context, id string) (*storage.Event, error)
	ListEventsFunc  func(ctx context.Context, userId string) ([]*storage.Event, error)
}

// common fixture for the Create/Update/... functions.
type TestFixture struct {
	T           *testing.T
	t0          *MockLogger
	mockStorage *MockStorage
	server      *Server
	handler     http.Handler
}

func (fx *TestFixture) LogBuffer() string {
	return fx.t0.buf.String()
}

// NewTestFixture создаёт базовый набор для тестов HTTP.
func NewTestFixture(t *testing.T, mockStorage *MockStorage) *TestFixture {
	t.Helper() // помечает функцию как вспомогательную для тестов.

	t0 := NewMockLogger()
	mockApp := &app.App{
		Logger: t0.logger,
	}

	server := NewServer(Config{
		Host: "localhost",
		Port: "8080",
	}, mockApp, mockStorage)

	server.RegisterHandlers()
	handler := server.GetHandler() // нужно добавить этот метод

	return &TestFixture{
		T:           t,
		t0:          t0,
		mockStorage: mockStorage,
		server:      server,
		handler:     handler,
	}
}

func (m *MockStorage) CreateEvent(ctx context.Context, event *storage.Event) error {
	if m.CreateEventFunc != nil {
		return m.CreateEventFunc(ctx, event)
	}
	return nil
}

func (m *MockStorage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	if m.UpdateEventFunc != nil {
		return m.UpdateEventFunc(ctx, event)
	}
	return nil
}

func (m *MockStorage) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	if m.GetEventFunc != nil {
		return m.GetEventFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockStorage) DeleteEvent(ctx context.Context, id string) error {
	if m.DeleteEventFunc != nil {
		return m.DeleteEventFunc(ctx, id)
	}
	return nil
}

func (m *MockStorage) ListEvents(ctx context.Context, userID string) ([]*storage.Event, error) {
	if m.ListEventsFunc != nil {
		return m.ListEventsFunc(ctx, userID)
	}
	return nil, nil
}

// CreateEvent ...
func TestCreateEvent_Success(t *testing.T) {
	mockStorage := &MockStorage{
		CreateEventFunc: func(_ context.Context, event *storage.Event) error {
			event.ID = "generated-id-456" // имитируем генерацию ID.
			return nil
		},
	}
	fx := NewTestFixture(t, mockStorage)

	// Подготовка запроса
	reqBody := CreateEventRequest{
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   time.Date(2026, 2, 25, 15, 30, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 2, 25, 16, 30, 0, 0, time.UTC),
		UserID:      "user-123",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(context.Background(), "POST",
		"/events", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем
	rr := httptest.NewRecorder()
	fx.handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp EventResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, "generated-id-456", resp.ID)
	assert.Equal(t, "Test Event", resp.Title)
	assert.Equal(t, "Test Description", resp.Description)
	assert.Equal(t, time.Date(2026, 2, 25, 15, 30, 0, 0, time.UTC), resp.StartTime)
	assert.Equal(t, time.Date(2026, 2, 25, 16, 30, 0, 0, time.UTC), resp.EndTime)
	assert.Equal(t, "user-123", resp.UserID)

	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] HTTP Create/Request: title="Test Event", user_id=user-123`))
	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] HTTP Create/Response: title="Test Event", user_id=user-123`))
	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] Request completed method: POST path: /events ip:  latency:`))

	/*
		logContent := fx.LogBuffer()
		_ = os.WriteFile("test_logs.txt", []byte(logContent), 0644)
	*/
}

func TestCreateEvent_JsonUnmarshallingError(t *testing.T) {
	fx := NewTestFixture(t, &MockStorage{})

	// Подготовка запроса
	reqBody := CreateEventRequest{
		Title:       "Test Event",
		Description: "Test Description",
		UserID:      "user-123",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	bodyBytes[0] = '[' // spoiling Json string replace first { at [
	req, _ := http.NewRequestWithContext(context.Background(), "POST",
		"/events", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем
	rr := httptest.NewRecorder()
	fx.handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] HTTP Create Error: invalid json`))
	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] Request completed method: POST path: /events ip:  latency:`))
}

func TestCreateEvent_EventCreatingFailure(t *testing.T) {
	expectedErr := fmt.Errorf("database connection failed")
	mockStorage := &MockStorage{
		CreateEventFunc: func(_ context.Context, _ *storage.Event) error {
			return expectedErr
		},
	}
	fx := NewTestFixture(t, mockStorage)

	// Подготовка запроса
	reqBody := CreateEventRequest{
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   time.Date(2026, 2, 25, 15, 30, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 2, 25, 16, 30, 0, 0, time.UTC),
		UserID:      "user-123",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(context.Background(), "POST",
		"/events", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем
	rr := httptest.NewRecorder()
	fx.handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)

	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] HTTP Create/Request: title="Test Event", user_id=user-123`))
	assert.True(t, strings.Contains(fx.LogBuffer(),
		`[I] HTTP Create Error: event creating failed database connection failed`))
	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] Request completed method: POST path: /events ip:  latency:`))
}

func Test_EventIDFromURL(t *testing.T) {
	testCases := []struct {
		name string
		url  string
		id   string
		errS bool
	}{
		{
			name: "правильный URL",
			url:  "/events/123",
			id:   "123",
			errS: false,
		},
		{
			name: "пустой ID",
			url:  "/events/",
			id:   "",
			errS: true,
		},
		{
			name: "нет ID",
			url:  "/events",
			id:   "",
			errS: true,
		},
		{
			name: "лишний сегмент",
			url:  "/events/123/extra",
			id:   "",
			errS: true,
		},
		{
			name: "короткий URL",
			url:  "/",
			id:   "",
			errS: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var s *Server
			id, err := s.EventIDFromURL(tc.url)

			assert.Equal(t, tc.id, id)
			if tc.errS {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// UpdateEvent ...
func TestUpdateEvent_Success(t *testing.T) {
	mockStorage := &MockStorage{
		UpdateEventFunc: func(_ context.Context, event *storage.Event) error {
			event.Description += " :updated"
			return nil
		},
	}
	fx := NewTestFixture(t, mockStorage)

	// Подготовка запроса
	ID := "generated-id-717"
	reqBody := UpdateEventRequest{
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   time.Date(2026, 2, 25, 15, 30, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 2, 25, 16, 30, 0, 0, time.UTC),
		UserID:      "user-123",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(context.Background(), "PUT",
		fmt.Sprintf("/events/%s", ID), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем
	rr := httptest.NewRecorder()
	fx.handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp EventResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, "generated-id-717", resp.ID)
	assert.Equal(t, "Test Event", resp.Title)
	assert.Equal(t, "Test Description :updated", resp.Description)
	assert.Equal(t, time.Date(2026, 2, 25, 15, 30, 0, 0, time.UTC), resp.StartTime)
	assert.Equal(t, time.Date(2026, 2, 25, 16, 30, 0, 0, time.UTC), resp.EndTime)
	assert.Equal(t, "user-123", resp.UserID)

	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] HTTP Update/Request: title="Test Event", user_id=user-123`))
	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] HTTP Update/Response: title="Test Event", user_id=user-123`))
	assert.True(t, strings.Contains(fx.LogBuffer(),
		`[I] Request completed method: PUT path: /events/generated-id-717 ip:  latency:`))
}

func TestUpdateEvent_JsonUnmarshallingError(t *testing.T) {
	fx := NewTestFixture(t, &MockStorage{})

	// Подготовка запроса
	ID := "generated-id-987"
	reqBody := UpdateEventRequest{
		Title:       "Test Event",
		Description: "Test Description",
		UserID:      "user-123",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	bodyBytes[0] = '[' // spoiling Json string replace first { at [
	req, _ := http.NewRequestWithContext(context.Background(), "PUT",
		fmt.Sprintf("/events/%s", ID), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем
	rr := httptest.NewRecorder()
	fx.handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] HTTP Update Error: invalid json`))
	assert.True(t, strings.Contains(fx.LogBuffer(),
		`[I] Request completed method: PUT path: /events/generated-id-987 ip:  latency:`))
}

func TestUpdateEvent_EventUpdatingFailure(t *testing.T) {
	expectedErr := fmt.Errorf("database connection failed")
	mockStorage := &MockStorage{
		UpdateEventFunc: func(_ context.Context, _ *storage.Event) error {
			return expectedErr
		},
	}
	fx := NewTestFixture(t, mockStorage)

	// Подготовка запроса
	ID := "generated-id-ABC"
	reqBody := UpdateEventRequest{
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   time.Date(2026, 2, 25, 15, 30, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 2, 25, 16, 30, 0, 0, time.UTC),
		UserID:      "user-123",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(context.Background(), "PUT",
		fmt.Sprintf("/events/%s", ID), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем
	rr := httptest.NewRecorder()
	fx.handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)

	assert.True(t, strings.Contains(fx.LogBuffer(), `[I] HTTP Update/Request: title="Test Event", user_id=user-123`))
	assert.True(t, strings.Contains(fx.LogBuffer(),
		`[I] HTTP Update Error: event updating failed database connection failed`))
	assert.True(t, strings.Contains(fx.LogBuffer(),
		`[I] Request completed method: PUT path: /events/generated-id-ABC ip:  latency:`))

	logContent := fx.LogBuffer()
	_ = os.WriteFile("test_logs.txt", []byte(logContent), 0644)
}
