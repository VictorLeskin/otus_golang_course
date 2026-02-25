// DTO (Data Transfer Objects) — объекты для передачи данных по сети.
// Они отделяют внешний JSON-формат от внутренних моделей (storage.Event).
package internalhttp

import "time"

// CreateEventRequest — DTO для создания события (POST /events)
// Клиент присылает JSON, который декодируется в эту структуру.
type CreateEventRequest struct {
	Title       string    `json:"title"`                 // Название события
	Description string    `json:"description,omitempty"` // Описание (необязательное)
	StartTime   time.Time `json:"start_time"`            // Время начала (RFC3339)
	EndTime     time.Time `json:"end_time"`              // Время окончания (RFC3339)
	UserID      string    `json:"user_id"`               // ID владельца
}

// UpdateEventRequest — DTO для обновления события (PUT /events/{id})
// Внимание: ID события передаётся в URL, не в теле!
type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	UserID      string    `json:"user_id"`
}

// EventResponse — DTO для ответа с данными события
// Возвращается при создании, получении и обновлении события.
type EventResponse struct {
	ID          string    `json:"id"` // Сгенерированный сервером ID
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	UserID      string    `json:"user_id"`
}

// ErrorResponse — DTO для ответа с ошибкой
// Все ошибки возвращаются в едином формате.
type ErrorResponse struct {
	Error string `json:"error"` // Текст ошибки
}

// ListEventsResponse — DTO для списка событий (GET /events)
type ListEventsResponse struct {
	Events []EventResponse `json:"events"` // Массив событий
}
