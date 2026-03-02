package app

import (
	"context"

	"calendar/internal/logger"
	"calendar/internal/storage"
)

type App struct { // TODO
	Logger  *logger.Logger
	Storage storage.Storage
}

func New(logger *logger.Logger, storage storage.Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	_ = ctx
	_ = id
	_ = title
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
