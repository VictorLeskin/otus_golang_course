package memorystorage

import (
	"calendar/internal/storage"
	"context"
	"sync"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	events map[string]*storage.Event
}

func New() *MemoryStorage {
	return &MemoryStorage{
		events: make(map[string]*storage.Event),
	}
}

func (ms *MemoryStorage) CreateEvent(ctx context.Context, event *storage.Event) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, exists := ms.events[event.ID]; exists {
		return storage.ErrEventExists
	}

	ms.events[event.ID] = event
	return nil
}

func (ms *MemoryStorage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, exists := ms.events[event.ID]; !exists {
		return storage.ErrEventNotFound
	}

	ms.events[event.ID] = event
	return nil
}

func (ms *MemoryStorage) DeleteEvent(ctx context.Context, id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, exists := ms.events[id]; !exists {
		return storage.ErrEventNotFound
	}

	delete(ms.events, id)

	return nil
}

func (ms *MemoryStorage) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	event, exists := ms.events[id]
	if !exists {
		return nil, storage.ErrEventNotFound
	}

	return event, nil
}

func (ms *MemoryStorage) ListEvents(ctx context.Context, userID string) ([]*storage.Event, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var result []*storage.Event
	for _, event := range ms.events {
		if event.UserID == userID {
			result = append(result, event)
		}
	}

	return result, nil
}

// TODO
