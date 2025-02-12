package memorystorage

import (
	"sync"
	"testing"

	//nolint:depguard
	"github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	e := storage.Event{
		UUID:           "3fdf95e9-59b8-499c-b249-41a67f0d2293",
		Summary:        "Заголовок нового события",
		StartedAt:      "2025-02-11 10:30:00",
		FinishedAt:     "2022-02-12 11:00:00",
		Description:    "Описание нового события",
		UserUUID:       "facecede-142d-4038-a9bd-05aa226a0449",
		NotificationAt: "2025-02-11 10:00:00",
	}
	s := MemoryStorage{
		mu:     sync.RWMutex{},
		events: map[string]storage.Event{},
	}

	// Создать объект события можно в случае, если событие ранее не было создано
	ok, err := s.CreateEvent(e)
	require.NoError(t, err)
	require.True(t, ok)
	require.Less(t, 0, len(s.events))

	// Создать событие ранее было создано, то при создании возникнет ошибка
	ok, err = s.CreateEvent(e)
	require.Error(t, err)
	require.False(t, ok)
	require.Equal(t, 1, len(s.events))
}

func TestUpdateEvent(t *testing.T) {
	uuid := "3fdf95e9-59b8-499c-b249-41a67f0d2293"
	event := storage.Event{
		UUID:           "3fdf95e9-59b8-499c-b249-41a67f0d2293",
		Summary:        "Заголовок нового события",
		StartedAt:      "2025-02-11 10:30:00",
		FinishedAt:     "2022-02-12 11:00:00",
		Description:    "Описание нового события",
		UserUUID:       "facecede-142d-4038-a9bd-05aa226a0449",
		NotificationAt: "2025-02-11 10:00:00",
	}

	s := MemoryStorage{
		mu:     sync.RWMutex{},
		events: map[string]storage.Event{},
	}

	// Попытка обновить ранее не созданное событие приведет к ошибке
	ok, err := s.UpdateEvent(uuid, event)
	require.Error(t, err)
	require.False(t, ok)
	ok, err = s.DeleteEvent(uuid)
	require.Error(t, err)
	require.False(t, ok)

	// Создать объект события можно в случае, если событие ранее не было создано
	ok, err = s.CreateEvent(event)
	require.NoError(t, err)
	require.True(t, ok)
	require.Less(t, 0, len(s.events))

	// Попытка обновить ранее существующее событие приведет пройдет без ошибок
	ok, err = s.UpdateEvent(uuid, event)
	require.NoError(t, err)
	require.True(t, ok)

	eventFound, err := s.GetEvent(uuid)
	require.NoError(t, err)
	require.Equal(t, uuid, eventFound.UUID)
	require.Equal(t, event.Summary, eventFound.Summary)
	require.Equal(t, event.StartedAt, eventFound.StartedAt)
	require.Equal(t, event.FinishedAt, eventFound.FinishedAt)
	require.Equal(t, event.Description, eventFound.Description)
	require.Equal(t, event.UserUUID, eventFound.UserUUID)
	require.Equal(t, event.NotificationAt, eventFound.NotificationAt)

	ok, err = s.DeleteEvent(uuid)
	require.NoError(t, err)
	require.True(t, ok)
	eventFound, err = s.GetEvent(uuid)
	require.Error(t, err)
}
