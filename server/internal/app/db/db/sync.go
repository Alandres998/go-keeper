package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Alandres998/go-keeper/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// AddSyncHistory добавляет запись о синхронизации данных
func (s *DBStorage) AddSyncHistory(ctx context.Context, tx *sqlx.Tx, history *models.SyncHistory) (*models.SyncHistory, error) {
	query := `
		INSERT INTO sync_history (user_id, sync_timestamp, data_version, operation_type, meta_info)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, sync_timestamp
	`

	now := time.Now()

	metaInfoJSON, err := json.Marshal(history.MetaInfo)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации MetaInfo: %v", err)
	}

	// Используем транзакцию, если она передана
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, history.UserID, now, history.DataVersion, history.OperationType, string(metaInfoJSON)).Scan(&history.ID, &history.SyncTimestamp)
	} else {
		err = s.db.QueryRowContext(ctx, query, history.UserID, now, history.DataVersion, history.OperationType, string(metaInfoJSON)).Scan(&history.ID, &history.SyncTimestamp)
	}

	if err != nil {
		return nil, fmt.Errorf("не смог загрузить историю: %v", err)
	}

	history.SyncTimestamp = now
	return history, nil
}
