package storage

type Event struct {
	// UUID - уникальный идентификатор события.
	UUID string `db:"uuid"`
	// Заголовок - короткий текст.
	Summary string `db:"summary"`
	// Дата и время начала события.
	StartedAt string `db:"started_at"`
	// Дата и время окончания события.
	FinishedAt string `db:"finished_at"`
	// Описание события.
	Description string `db:"description"`
	// UUID пользователя, владельца события.
	UserUUID string `db:"user_uuid"`
	// Дата и время уведомления о событии.
	NotificationAt string `db:"notification_at"`
}
