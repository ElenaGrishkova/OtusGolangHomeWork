package sqlstorage

import (
	//nolint:depguard
	"github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/app"
	//nolint:depguard
	"github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/jmoiron/sqlx"
	// Необходимо импортировать пакет для того чтобы подключился драйвер pq.
	//nolint:depguard
	_ "github.com/lib/pq"
)

var (
	schema = `
-- Схема для хранения событий в базе данных
CREATE SCHEMA IF NOT EXISTS "events";

-- Таблица для хранения событий в базе данных
CREATE TABLE IF NOT EXISTS "events"."events"
(
	-- Get - уникальный идентификатор события
    "uuid" varchar
		constraint events_pk
			primary key,
	-- Заголовок - короткий текст
	"summary" varchar not null,
	-- Дата и время начала события
	"started_at" varchar not null,
	-- Дата и время начала события
	"finished_at" varchar not null,
	-- Описание события
	"description" varchar not null,
	-- Get пользователя, владельца события
	"user_uuid" varchar not null,
	-- Дата и время уведомления о событии
	"notification_at" varchar not null
);`

	sqlEventSelectByID = `SELECT * FROM "events"."events" WHERE "uuid" = $1 LIMIT 1`

	sqlEventInsert = `-- Запрос создающий запись в базе данных о событии 
INSERT INTO "events"."events"
(uuid, summary, started_at, finished_at, description, user_uuid, notification_at)
VALUES (:uuid, :summary, :started_at, :finished_at, :description, :user_uuid, :notification_at)`

	sqlEventUpdate = `-- Запрос обновляющий запись в базе данных о событии
UPDATE "events"."events"
SET summary = $2,
    started_at = $3,
    finished_at = $4,
    description = $5,
    user_uuid = $6,
    notification_at = $7
WHERE uuid = $1`

	sqlEventDelete = `-- Запрос удаляющий запись в базе данных о событии
DELETE "events"."events"
WHERE uuid = $1`
)

type SQLStorage struct {
	driver    string
	dsn       string
	dbConnect *sqlx.DB
}

func NewSQLStorage(driver string, dsn string) app.Storage {
	return &SQLStorage{
		driver: driver,
		dsn:    dsn,
	}
}

// Добавляет в базу данных новое событие.
func (s *SQLStorage) CreateEvent(event storage.Event) (bool, error) {
	tx := s.dbConnect.MustBegin()
	if _, err := tx.NamedExec(sqlEventInsert, &event); err != nil {
		_ = tx.Rollback()
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

// Обновляет существующее в хранилище событие.
func (s *SQLStorage) UpdateEvent(uuid string, attributes storage.Event) (bool, error) {
	tx := s.dbConnect.MustBegin()
	args := []interface{}{
		uuid,
		attributes.Summary,
		attributes.StartedAt,
		attributes.FinishedAt,
		attributes.Description,
		attributes.UserUUID,
		attributes.NotificationAt,
	}
	if _, err := tx.Exec(sqlEventUpdate, args...); err != nil {
		_ = tx.Rollback()
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

// Удаляет существующее событие из хранилища.
func (s *SQLStorage) DeleteEvent(uuid string) (bool, error) {
	tx := s.dbConnect.MustBegin()
	args := []interface{}{
		uuid,
	}
	if _, err := tx.Exec(sqlEventDelete, args...); err != nil {
		_ = tx.Rollback()
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

// Возвращает список соответствующих условию событий из хранилища, проиндексированные по идентификатору.
func (s *SQLStorage) SelectEvents() (map[string]storage.Event, error) {
	events := map[string]storage.Event{}
	items := []storage.Event{}
	err := s.dbConnect.Select(&items, `SELECT * FROM "events"."events"`)
	if err != nil {
		return events, err
	}
	for _, event := range items {
		events[event.UUID] = event
	}
	return events, nil
}

// Возвращает событие из хранилища по идентификатору.
func (s *SQLStorage) GetEvent(uuid string) (storage.Event, error) {
	event := storage.Event{}
	rows, err := s.dbConnect.Queryx(sqlEventSelectByID, uuid)
	if err != nil {
		return storage.Event{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.StructScan(&event)
		if err != nil {
			return storage.Event{}, err
		}
	}
	return event, nil
}

func (s *SQLStorage) Connect() error {
	db, err := sqlx.Connect(s.driver, s.dsn)
	if err == nil {
		// Миграция
		db.MustExec(schema)
		s.dbConnect = db
	}

	return err
}

func (s *SQLStorage) Close() error {
	if s.dbConnect != nil {
		return s.dbConnect.Close()
	}

	return nil
}
