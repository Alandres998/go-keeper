package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Alandres998/go-keeper/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// InsertOrUpdatePrivateData добавляет или обновляет запись в таблице private_data.
func (s *DBStorage) InsertOrUpdatePrivateData(ctx context.Context, tx *sqlx.Tx, data *models.PrivateData) (*models.PrivateData, error) {
	query := `
		INSERT INTO private_data (user_id, text_data, binary_data, card_number, created_at, updated_at, meta_info)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), $5)
		ON CONFLICT (user_id)
		DO UPDATE SET
			text_data = EXCLUDED.text_data,
			binary_data = EXCLUDED.binary_data,
			card_number = EXCLUDED.card_number,
			meta_info = EXCLUDED.meta_info,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	metaInfoJSON, err := json.Marshal(data.MetaInfo)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации meta_info: %v", err)
	}

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, data.UserID, data.TextData, data.BinaryData, data.CardNumber, metaInfoJSON)
	} else {
		row = s.db.QueryRowContext(ctx, query, data.UserID, data.TextData, data.BinaryData, data.CardNumber, metaInfoJSON)
	}

	err = row.Scan(&data.ID, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}

	return data, nil
}
