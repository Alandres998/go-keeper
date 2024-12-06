package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Alandres998/go-keeper/server/internal/app/db/storage"
	"github.com/Alandres998/go-keeper/server/internal/app/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
		data_type VARCHAR(50) NOT NULL,
		data_key VARCHAR(255) NOT NULL,
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

// CreateUser создает нового пользователя в базе данных
func (s *DBStorage) CreateUser(ctx context.Context, tx *sqlx.Tx, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (login, password, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	now := time.Now()
	var err error

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования пароля: %v", err)
	}
	user.Password = string(hashedPassword)
	// Используем транзакцию, если она передана
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, user.Login, user.Password, user.Email, now, now).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	} else {
		err = s.db.QueryRowContext(ctx, query, user.Login, user.Password, user.Email, now, now).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("не удалось создать пользователя, никаких строк не вернулось: %v", err)
		}
		return nil, fmt.Errorf("не смог создать пользователя: %v", err)
	}

	return user, nil
}

// GetUserByLogin получает пользователя по логину
func (s *DBStorage) GetUserByLogin(ctx context.Context, tx *sqlx.Tx, login string) (*models.User, error) {
	query := `
		SELECT id, login, password, email, created_at, updated_at
		FROM users
		WHERE login = $1
	`
	var user models.User
	var err error

	// Используем транзакцию, если она передана
	if tx != nil {
		err = tx.GetContext(ctx, &user, query, login)
	} else {
		err = s.db.GetContext(ctx, &user, query, login)
	}

	if err != nil {
		return nil, fmt.Errorf("не корректный логин или пароль: %v", err)
	}

	return &user, nil
}

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

// CountUserLogins возвращает количество попыток входа для пользователя
func (s *DBStorage) CountUserLogins(ctx context.Context, tx *sqlx.Tx, userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM sync_history WHERE user_id = $1 AND operation_type = 'login'`

	// Используем транзакцию, если она передана
	if tx != nil {
		err := tx.QueryRowContext(ctx, query, userID).Scan(&count)
		if err != nil {
			return 0, err
		}
	} else {
		err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

// InsertPrivateData добавляет запись в таблицу private_data.
func (s *DBStorage) InsertPrivateData(ctx context.Context, tx *sqlx.Tx, data *models.PrivateData) (*models.PrivateData, error) {
	query := `
		INSERT INTO private_data (user_id, data_type, data_key, text_data, binary_data, card_number, created_at, updated_at, meta_info)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), $7)
		RETURNING id, created_at, updated_at
	`

	metaInfoJSON, err := json.Marshal(data.MetaInfo)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации meta_info: %v", err)
	}

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, data.UserID, data.DataType, data.DataKey, data.TextData, data.BinaryData, data.CardNumber, metaInfoJSON)
	} else {
		row = s.db.QueryRowContext(ctx, query, data.UserID, data.DataType, data.DataKey, data.TextData, data.BinaryData, data.CardNumber, metaInfoJSON)
	}

	err = row.Scan(&data.ID, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}

	return data, nil
}

// GetPrivateDataByUserID возвращает данные пользователя по userID
func (s *DBStorage) GetPrivateDataByUserID(ctx context.Context, userID int) ([]*models.PrivateData, error) {
	query := `
		SELECT id, user_id, card_number, text_data, binary_data, meta_info
		FROM private_data
		WHERE user_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var privateDataList []*models.PrivateData

	for rows.Next() {
		var data models.PrivateData
		if err := rows.Scan(&data.ID, &data.UserID, &data.CardNumber, &data.TextData, &data.BinaryData, &data.MetaInfo); err != nil {
			return nil, err
		}
		privateDataList = append(privateDataList, &data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return privateDataList, nil
}
