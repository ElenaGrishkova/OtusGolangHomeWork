package app

import (
	"context"

	//nolint:depguard
	"github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	Connect() error
	Close() error
	CreateEvent(event storage.Event) (bool, error)
	UpdateEvent(uuid string, event storage.Event) (bool, error)
	DeleteEvent(uuid string) (bool, error)
	SelectEvents() (map[string]storage.Event, error)
	GetEvent(uuid string) (storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{storage: storage, logger: logger}
}

func (a *App) CreateEvent(_ context.Context, id, title string) error {
	ok, err := a.storage.CreateEvent(storage.Event{UUID: id, Summary: title})
	if !ok {
		return err
	}
	return nil
}

// TODO
