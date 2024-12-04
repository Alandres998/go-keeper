package storage

import (
	"context"

	"github.com/Alandres998/go-keeper/server/internal/app/models"
	"github.com/jmoiron/sqlx"
)

var Store Storage

// Storage интерфейс определяет методы для работы с ним
type Storage interface {
	// Начало транзакции
	BeginTx(ctx context.Context) (*sqlx.Tx, error)

	// Методы, работающие с транзакциями
	CreateUser(ctx context.Context, tx *sqlx.Tx, user *models.User) (*models.User, error)
	GetUserByLogin(ctx context.Context, tx *sqlx.Tx, login string) (*models.User, error)
	AddSyncHistory(ctx context.Context, tx *sqlx.Tx, history *models.SyncHistory) (*models.SyncHistory, error)
	CountUserLogins(ctx context.Context, tx *sqlx.Tx, userID int) (int, error)
}
