package grpc

import (
	"calendar/api/pb/calendar"
	"calendar/internal/logger"
	"calendar/internal/storage"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// internal/server/grpc/mock_storage_test.go
type MockStorage struct {
	CreateEventFunc func(ctx context.Context, event *storage.Event) error
	UpdateEventFunc func(ctx context.Context, event *storage.Event) error
	DeleteEventFunc func(ctx context.Context, id string) error
	// остальные методы можно добавить по необходимости
}

// MockLogger обёртка над logger.NewWriterLogger для тестов
type MockLogger struct {
	buf    strings.Builder
	logger *logger.Logger
}

func NewMockLogger() *MockLogger {
	ml := &MockLogger{}
	ml.logger = logger.NewWriterLogger("info", &ml.buf)
	return ml
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

// Заглушки для остальных методов интерфейса (чтобы MockStorage удовлетворял интерфейсу)
func (m *MockStorage) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	return nil, nil
}
func (m *MockStorage) DeleteEvent(ctx context.Context, id string) error {
	if m.DeleteEventFunc != nil {
		return m.DeleteEventFunc(ctx, id)
	}
	return nil
}
func (m *MockStorage) ListEvents(ctx context.Context, userID string) ([]*storage.Event, error) {
	return nil, nil
}

func TestConvertFromPBEvent(t *testing.T) {
	// Создаем тестовые данные
	fixedTime := time.Date(2026, time.February, 23, 15, 30, 0, 0, time.UTC)

	pbEvent := &calendar.Event{
		Id:          "123",
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   timestamppb.New(fixedTime),
		EndTime:     timestamppb.New(fixedTime.Add(30 * time.Minute)),
		UserId:      "user-456",
	}

	// Вызываем тестируемую функцию
	storageEvent := convertFromPBEvent(pbEvent)

	// Проверяем все поля
	assert.Equal(t, "123", storageEvent.ID)
	assert.Equal(t, "Test Event", storageEvent.Title)
	assert.Equal(t, "Test Description", storageEvent.Description)
	assert.Equal(t, time.Date(2026, time.February, 23, 15, 30, 0, 0, time.UTC), storageEvent.StartTime)
	assert.Equal(t, time.Date(2026, time.February, 23, 16, 0, 0, 0, time.UTC), storageEvent.EndTime)
	assert.Equal(t, "user-456", storageEvent.UserID)
}

func TestConvertToPBEvent(t *testing.T) {
	// Создаем тестовые данные с фиксированными значениями
	fixedTime := time.Date(2026, time.February, 23, 15, 30, 0, 0, time.UTC)

	storageEvent := &storage.Event{
		ID:          "event-789",
		Title:       "Daily Standup",
		Description: "15-minute team sync",
		StartTime:   fixedTime,
		EndTime:     fixedTime.Add(30 * time.Minute),
		UserID:      "user-101",
	}

	// Вызываем тестируемую функцию
	pbEvent := convertToPBEvent(storageEvent)

	// Проверяем все поля через assert с КОНКРЕТНЫМИ значениями
	assert.Equal(t, "event-789", pbEvent.Id)
	assert.Equal(t, "Daily Standup", pbEvent.Title)
	assert.Equal(t, "15-minute team sync", pbEvent.Description)
	assert.Equal(t, time.Date(2026, time.February, 23, 15, 30, 0, 0, time.UTC), pbEvent.StartTime.AsTime())
	assert.Equal(t, time.Date(2026, time.February, 23, 16, 0, 0, 0, time.UTC), pbEvent.EndTime.AsTime())
	assert.Equal(t, "user-101", pbEvent.UserId)
}

func TestLogError(t *testing.T) {
	t0 := NewMockLogger()

	server := &Server{
		logger: t0.logger,
	}

	// Создаем тестовую ошибку
	testErr := fmt.Errorf("database connection failed")

	// Вызываем тестируемую функцию
	server.LogError("Create", testErr)

	// Проверяем, что логгер получил сообщение
	assert.Equal(t, "[I] gRPC Create Error: database connection failed\n", t0.buf.String())
}

func TestLogCalendarEvent(t *testing.T) {
	t0 := NewMockLogger()

	server := &Server{
		logger: t0.logger,
	}

	// Создаем тестовое событие
	event := &calendar.Event{
		Id:          "event-789",
		Title:       "Daily Standup",
		Description: "15-minute team sync",
		StartTime:   timestamppb.New(time.Date(2026, time.February, 23, 15, 30, 0, 0, time.UTC)),
		EndTime:     timestamppb.New(time.Date(2026, time.February, 23, 16, 0, 0, 0, time.UTC)),
		UserId:      "user-101",
	}

	// Вызываем тестируемую функцию
	server.LogCalendarEvent("Create", "Request", event)

	// Проверяем, что логгер получил сообщение
	assert.True(t, strings.HasPrefix(t0.buf.String(), "[I] gRPC Create/Request"))

	expectedParts := []string{
		"Id: event-789",
		`Title: "Daily Standup"`,
		`Description: "15-minute team sync"`,
		"StartTime:2026-02-23T15:30:00Z",
		"EndTime:2026-02-23T16:00:00Z",
		"UserId: user-101",
	}

	for _, part := range expectedParts {
		assert.Contains(t, t0.buf.String(), part)
	}
}

func TestLogDeleteRequest(t *testing.T) {
	t0 := NewMockLogger()

	server := &Server{
		logger: t0.logger,
	}

	req := &calendar.DeleteEventRequest{
		Id: "id1",
	}
	server.LogDeleteRequest(req)

	assert.Equal(t, "[I] gRPC Delete/Request Id: id1\n", t0.buf.String())
}

func TestLogDeleteResponse(t *testing.T) {
	t0 := NewMockLogger()

	server := &Server{
		logger: t0.logger,
	}

	req := &calendar.DeleteEventRequest{
		Id: "id1",
	}
	server.LogDeleteResponse(req)

	assert.Equal(t, "[I] gRPC Delete/Response Id: id1\n", t0.buf.String())
}

func TestCreateEvent_Success(t *testing.T) {
	t0 := NewMockLogger()

	mockStorage := &MockStorage{
		CreateEventFunc: func(ctx context.Context, event *storage.Event) error {
			event.ID = "generated-id-123" // имитируем генерацию ID
			return nil
		},
	}

	server := &Server{
		storage: mockStorage,
		logger:  t0.logger,
	}

	req := &calendar.CreateEventRequest{
		Event: &calendar.Event{
			Title:       "Test Event",
			Description: "Test Description",
			StartTime:   timestamppb.New(time.Date(2026, 2, 23, 15, 30, 0, 0, time.UTC)),
			EndTime:     timestamppb.New(time.Date(2026, 2, 23, 16, 0, 0, 0, time.UTC)),
			UserId:      "user-123",
		},
	}

	resp, err := server.CreateEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Event)
	assert.Equal(t, "generated-id-123", resp.Event.Id)
	assert.Equal(t, "Test Event", resp.Event.Title)
	assert.Equal(t, "Test Description", resp.Event.Description)
	assert.Equal(t, time.Date(2026, time.February, 23, 15, 30, 0, 0, time.UTC), resp.Event.StartTime.AsTime())
	assert.Equal(t, time.Date(2026, time.February, 23, 16, 0, 0, 0, time.UTC), resp.Event.EndTime.AsTime())
	assert.Equal(t, "user-123", resp.Event.UserId)

	// Проверка логов
	assert.True(t, strings.HasPrefix(t0.buf.String(), "[I] gRPC Create/Request"))
	assert.True(t, strings.Contains(t0.buf.String(), "[I] gRPC Create/Response"))
}

func TestCreateEvent_StorageError(t *testing.T) {
	t0 := NewMockLogger()

	expectedErr := fmt.Errorf("database connection failed")
	mockStorage := &MockStorage{
		CreateEventFunc: func(ctx context.Context, event *storage.Event) error {
			return expectedErr
		},
	}

	server := &Server{
		storage: mockStorage,
		logger:  t0.logger,
	}

	// Подготовка запроса
	req := &calendar.CreateEventRequest{
		Event: &calendar.Event{
			Title:       "Test Event",
			Description: "Test Description",
			StartTime:   timestamppb.New(time.Date(2026, 2, 23, 15, 30, 0, 0, time.UTC)),
			EndTime:     timestamppb.New(time.Date(2026, 2, 23, 16, 0, 0, 0, time.UTC)),
			UserId:      "user-123",
		},
	}

	// Вызов
	resp, err := server.CreateEvent(context.Background(), req)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NotNil(t, resp)
	assert.Nil(t, resp.Event)
	assert.Equal(t, expectedErr.Error(), resp.ErrorMessage)

	// Проверка логов
	logOutput := t0.buf.String()
	assert.Contains(t, logOutput, "Create/Request")
	assert.Contains(t, logOutput, "Create Error")
	assert.Contains(t, logOutput, "database connection failed")
	assert.NotContains(t, logOutput, "Create/Response")
}

// UpdateEvent
func TestUpdateEvent_Success(t *testing.T) {
	t0 := NewMockLogger()

	mockStorage := &MockStorage{
		UpdateEventFunc: func(ctx context.Context, event *storage.Event) error {
			event.Description += " :updated"
			return nil
		},
	}

	server := &Server{
		storage: mockStorage,
		logger:  t0.logger,
	}

	req := &calendar.UpdateEventRequest{
		Event: &calendar.Event{
			Id:          "id-12222",
			Title:       "Test Event",
			Description: "Test Description",
			StartTime:   timestamppb.New(time.Date(2026, 2, 23, 15, 30, 0, 0, time.UTC)),
			EndTime:     timestamppb.New(time.Date(2026, 2, 23, 16, 0, 0, 0, time.UTC)),
			UserId:      "user-123",
		},
	}

	resp, err := server.UpdateEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Event)
	assert.Equal(t, "id-12222", resp.Event.Id)
	assert.Equal(t, "Test Event", resp.Event.Title)
	assert.Equal(t, "Test Description :updated", resp.Event.Description)
	assert.Equal(t, time.Date(2026, time.February, 23, 15, 30, 0, 0, time.UTC), resp.Event.StartTime.AsTime())
	assert.Equal(t, time.Date(2026, time.February, 23, 16, 0, 0, 0, time.UTC), resp.Event.EndTime.AsTime())
	assert.Equal(t, "user-123", resp.Event.UserId)

	// Проверка логов
	assert.True(t, strings.HasPrefix(t0.buf.String(), "[I] gRPC Update/Request"))
	assert.True(t, strings.Contains(t0.buf.String(), "[I] gRPC Update/Response"))
}

func TestUpdateEvent_StorageError(t *testing.T) {
	t0 := NewMockLogger()

	expectedErr := fmt.Errorf("database connection failed")
	mockStorage := &MockStorage{
		UpdateEventFunc: func(ctx context.Context, event *storage.Event) error {
			return expectedErr
		},
	}

	server := &Server{
		storage: mockStorage,
		logger:  t0.logger,
	}

	// Подготовка запроса
	req := &calendar.UpdateEventRequest{
		Event: &calendar.Event{
			Title:       "Test Event",
			Description: "Test Description",
			StartTime:   timestamppb.New(time.Date(2026, 2, 23, 15, 30, 0, 0, time.UTC)),
			EndTime:     timestamppb.New(time.Date(2026, 2, 23, 16, 0, 0, 0, time.UTC)),
			UserId:      "user-123",
		},
	}

	// Вызов
	resp, err := server.UpdateEvent(context.Background(), req)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NotNil(t, resp)
	assert.Nil(t, resp.Event)
	assert.Equal(t, expectedErr.Error(), resp.ErrorMessage)

	// Проверка логов
	logOutput := t0.buf.String()
	assert.Contains(t, logOutput, "Update/Request")
	assert.Contains(t, logOutput, "Update Error")
	assert.Contains(t, logOutput, "database connection failed")
	assert.NotContains(t, logOutput, "Update/Response")
}

// DeleteEvent
func TestDeleteEvent_Success(t *testing.T) {
	t0 := NewMockLogger()

	mockStorage := &MockStorage{
		DeleteEventFunc: func(ctx context.Context, id string) error {
			_ = ctx
			_ = id
			return nil
		},
	}

	server := &Server{
		storage: mockStorage,
		logger:  t0.logger,
	}

	req := &calendar.DeleteEventRequest{
		Id: "id-225",
	}

	resp, err := server.DeleteEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Проверка логов
	assert.Equal(t, "[I] gRPC Delete/Request Id: id-225\n"+
		"[I] gRPC Delete/Response Id: id-225\n", t0.buf.String())
}

func TestDeleteEvent_StorageError(t *testing.T) {
	t0 := NewMockLogger()

	expectedErr := fmt.Errorf("database connection failed")
	mockStorage := &MockStorage{
		DeleteEventFunc: func(ctx context.Context, id string) error {
			_ = ctx
			_ = id
			return expectedErr
		},
	}

	server := &Server{
		storage: mockStorage,
		logger:  t0.logger,
	}

	req := &calendar.DeleteEventRequest{
		Id: "id-225",
	}

	resp, err := server.DeleteEvent(context.Background(), req)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.ErrorMessage)
	assert.Equal(t, expectedErr.Error(), resp.ErrorMessage)

	// Проверка логов
	assert.Equal(t, "[I] gRPC Delete/Request Id: id-225\n"+
		"[I] gRPC Delete Error: database connection failed\n", t0.buf.String())
}
