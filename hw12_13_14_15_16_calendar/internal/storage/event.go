package storage

import (
	"context"
	"fmt"
	"time"
)

type Event struct {
	ID          string
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	UserID      string
}

type Storage interface {
	CreateEvent(ctx context.Context, event *Event) error
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (*Event, error)
	ListEvents(ctx context.Context, userID string) ([]*Event, error)
}

var (
	ErrEventNotFound = fmt.Errorf("an event not found")
	ErrEventExists   = fmt.Errorf("this event exists")
)
