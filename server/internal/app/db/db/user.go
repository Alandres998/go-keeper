package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Alandres998/go-keeper/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

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

// GetPrivateDataByUserID возвращает данные пользователя по userID
func (s *DBStorage) GetPrivateDataByUserID(ctx context.Context, userID int) (*models.PrivateData, error) {
	query := `
		SELECT id, user_id, card_number, text_data, binary_data, meta_info
		FROM private_data
		WHERE user_id = $1
		LIMIT 1
	`

	var data models.PrivateData
	row := s.db.QueryRowContext(ctx, query, userID)

	// Сканируем полученные данные в структуру
	err := row.Scan(&data.ID, &data.UserID, &data.CardNumber, &data.TextData, &data.BinaryData, &data.MetaInfo)
	if err == sql.ErrNoRows {
		// Если данных нет, возвращаем nil и ошибку
		return nil, nil
	} else if err != nil {
		// Обрабатываем другие ошибки
		return nil, err
	}

	return &data, nil
}
