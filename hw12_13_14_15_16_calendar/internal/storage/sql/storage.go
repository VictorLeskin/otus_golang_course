package sqlstorage

import (
	"calendar/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq" // драйвер PostgreSQL
)

type Config struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	SSLMode  string
}

type SQLStorage struct { // TODO
	cfg Config
	db  *sql.DB
}

func New(cfg Config) *SQLStorage {
	return &SQLStorage{
		cfg: cfg,
	}
}

func (s *SQLStorage) DSN() string {
	c := s.cfg
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}

func (s *SQLStorage) Connect(ctx context.Context) error {
	db, err := sql.Open("postgres", s.DSN())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// check the connection
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	s.db = db

	return nil
}

func (s *SQLStorage) Close(_ context.Context) error {
	return s.db.Close()
}

// Вспомогательная функция для определения duplicate key.
func isDuplicateError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" { // PostgreSQL error code for duplicating "23505".
			return true
		}
	}

	return false
}

func (s *SQLStorage) CreateEvent(ctx context.Context, event *storage.Event) error {
	query := `
        INSERT INTO events (id, title, description, start_time, end_time, user_id)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := s.db.ExecContext(ctx, query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.UserID,
	)
	if err != nil {
		// Проверяем на duplicate key
		if isDuplicateError(err) {
			return storage.ErrEventExists
		}
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (s *SQLStorage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	query := `
        UPDATE events 
        SET title = $2, description = $3, start_time = $4, end_time = $5, user_id = $6
        WHERE id = $1
    `

	result, err := s.db.ExecContext(ctx, query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update event %s: %w", event.ID, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *SQLStorage) DeleteEvent(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event %s: %w", id, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *SQLStorage) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	query := `
        SELECT id, title, description, start_time, end_time, user_id
        FROM events WHERE id = $1
    `

	var event storage.Event
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.UserID,
	)

	if err == sql.ErrNoRows {
		return nil, storage.ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get event %s: %w", id, err)
	}

	return &event, nil
}

func (s *SQLStorage) ListEvents(ctx context.Context, userID string) ([]*storage.Event, error) {
	query := `
        SELECT id, title, description, start_time, end_time, user_id
        FROM events
        WHERE user_id = $1 
        ORDER BY start_time
    `

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list events for user %s: %w", userID, err)
	}
	defer rows.Close()

	var events []*storage.Event
	for rows.Next() {
		var event storage.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return events, nil
}
