package db

import (
	"context"
	"time"

	"github.com/Alandres998/go-keeper/server/internal/app/db/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// DBStorage для работы с базой данных.
type DBStorage struct {
	db *sqlx.DB
}

// Инициализация хранилища
func NewDBStorage(dsn string) (storage.Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = logger.Sync()
	}()

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("Проблемы при подключении к БД",
			zap.String("Не смог подключиться к БД", err.Error()),
		)
		return nil, err
	}

	createTableQuery := `
	    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        login VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );
	
	CREATE TABLE IF NOT EXISTS private_data (
		id SERIAL PRIMARY KEY, 
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		text_data TEXT,
		binary_data BYTEA,
		card_number VARCHAR(19),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		meta_info JSONB
	);

	CREATE TABLE IF NOT EXISTS sync_history (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE, 
		sync_timestamp TIMESTAMPTZ DEFAULT NOW(), 
		data_version INT,
		operation_type VARCHAR(50),
		meta_info JSONB
	);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, createTableQuery)
	if err != nil {
		logger.Error("Не удалось создать таблицу",
			zap.String("Ошибка", err.Error()),
		)
		return nil, err
	}

	return &DBStorage{db: db}, nil
}

// BeginTx начинает новую транзакцию
func (s *DBStorage) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return s.db.BeginTxx(ctx, nil)
}

// execContext проверяет, передана ли транзакция, и возвращает соответствующее подключение
func (s *DBStorage) execContext(ctx context.Context, tx *sqlx.Tx) sqlx.ExtContext {
	if tx != nil {
		return tx
	}
	return s.db
}
