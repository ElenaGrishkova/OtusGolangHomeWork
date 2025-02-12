package memorystorage

import (
	"errors"
	"sync"

	//nolint:depguard
	"github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/app"
	//nolint:depguard
	"github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrFailCreateEvent = errors.New("error creating event")
	ErrEventNotExists  = errors.New("error event not exists by id")
)

type MemoryStorage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func NewMemoryStorage() app.Storage {
	return &MemoryStorage{
		mu:     sync.RWMutex{},
		events: map[string]storage.Event{},
	}
}

// Создает новое событие.
func (s *MemoryStorage) CreateEvent(e storage.Event) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[e.UUID]; ok {
		return false, ErrFailCreateEvent
	}
	s.events[e.UUID] = e

	return true, nil
}

// Обновляет существующее в хранилище событие.
func (s *MemoryStorage) UpdateEvent(uuid string, event storage.Event) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[uuid]
	if !ok {
		return false, ErrEventNotExists
	}
	event.UUID = uuid

	s.events[uuid] = event

	return true, nil
}

// Удаляет существующее событие из хранилища.
func (s *MemoryStorage) DeleteEvent(uuid string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.events[uuid]
	if !ok {
		return false, ErrEventNotExists
	}
	delete(s.events, uuid)
	return true, nil
}

// Возвращает список соответствующих условию событий из хранилища, проиндексированные по идентификатору.
func (s *MemoryStorage) SelectEvents() (map[string]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.events, nil
}

// Возвращает событие из хранилища по идентификатору.
func (s *MemoryStorage) GetEvent(uuid string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var event storage.Event
	var ok bool

	if event, ok = s.events[uuid]; ok {
		return event, nil
	}

	return event, ErrEventNotExists
}

func (s *MemoryStorage) Connect() error {
	return nil
}

func (s *MemoryStorage) Close() error {
	return nil
}
