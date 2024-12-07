package models

import "time"

// SyncHistory представляет собой структуру для хранения информации о синхронизации данных.
type SyncHistory struct {
	ID            int       `db:"id"`             // Уникальный идентификатор записи синхронизации
	UserID        int       `db:"user_id"`        // Внешний ключ на пользователя
	SyncTimestamp time.Time `db:"sync_timestamp"` // Время синхронизации
	DataVersion   int       `db:"data_version"`   // Версия данных, которая синхронизирована
	OperationType string    `db:"operation_type"` // Тип операции (insert, update, delete)
	MetaInfo      string    `db:"meta_info"`      // Метаинформация о синхронизации (например, причина синхронизации)
}
