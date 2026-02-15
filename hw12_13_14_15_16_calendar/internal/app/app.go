package app

import (
	"calendar/internal/logger"
	"calendar/internal/storage"
	"context"
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
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
