package grpc

import (
	"calendar/api/pb/calendar"
	"calendar/internal/logger"
	"calendar/internal/storage"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
	// Создаем мок логгера и сервер
	// Создаем мок логгера и сервер
	var buf strings.Builder // уже имеет Write([]byte) (int, error)
	t0 := logger.NewWriterLogger("info", &buf)

	server := &Server{
		logger: t0,
	}

	// Создаем тестовую ошибку
	testErr := fmt.Errorf("database connection failed")

	// Вызываем тестируемую функцию
	server.LogError("Create", testErr)

	// Проверяем, что логгер получил сообщение
	assert.Equal(t, "[I] gRPC Create Error: database connection failed\n", buf.String())
}

func TestLogCalendarEvent(t *testing.T) {
	// Создаем мок логгера и сервер
	var buf strings.Builder // СѓР¶Рµ РёРјРµРµС‚ Write([]byte) (int, error)
	t0 := logger.NewWriterLogger("info", &buf)

	server := &Server{
		logger: t0,
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
	assert.True(t, strings.HasPrefix(buf.String(), "[I] gRPC"))

	expectedParts := []string{
		"Create/Request",
		"Id: event-789",
		`Title: "Daily Standup"`,
		`Description: "15-minute team sync"`,
		"StartTime:2026-02-23T15:30:00Z",
		"EndTime:2026-02-23T16:00:00Z",
		"UserId: user-101",
	}

	for _, part := range expectedParts {
		assert.Contains(t, buf.String(), part)
	}
}
