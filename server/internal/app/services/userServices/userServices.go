package userServices

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Alandres998/go-keeper/models"
	authlogin "github.com/Alandres998/go-keeper/server/internal/app/auth"
	"github.com/Alandres998/go-keeper/server/internal/app/db/storage"
	configserver "github.com/Alandres998/go-keeper/server/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// UserCreate создать пользователя в системе
func UserCreate(ctx context.Context, login string, passowrd string, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, configserver.Options.TimeOut)
	defer cancel()
	// Начинаем транзакцию
	tx, err := storage.Store.BeginTx(ctx)
	if err != nil {
		return nil, errors.New("Ошибка при создании транзакции")
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Создаем пользователя
	user := models.User{
		Login:    login,
		Password: passowrd,
		Email:    email,
	}

	// Используем транзакцию при создании пользователя
	createdUser, err := storage.Store.CreateUser(ctx, tx, &user)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Ошибка при регистрации пользователя: %v", err))
	}

	// Добавление истории синхронизации
	history := models.SyncHistory{
		UserID:        createdUser.ID,
		DataVersion:   1,
		OperationType: "registration",
		MetaInfo:      "Регистрация пользователя",
	}

	_, err = storage.Store.AddSyncHistory(ctx, tx, &history)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Ошибки при добавлении истории: %v", err))
	}

	// Коммит транзакции
	if err = tx.Commit(); err != nil {
		return nil, errors.New(fmt.Sprintf("Ошибка фиксации транзакции: %v", err))
	}
	return createdUser, nil
}

func UserLogin(ctx context.Context, login string, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, configserver.Options.TimeOut)
	defer cancel()

	tx, err := storage.Store.BeginTx(ctx)
	if err != nil {
		return nil, errors.New("Ошибка при создании транзакции")
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	user, err := storage.Store.GetUserByLogin(ctx, tx, login)
	if err != nil {
		return nil, errors.New("Не корректный логин или пароль")
	}

	// Сравниваем хешированный пароль

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("Не корректный логин или пароль")
	}

	dataVersion, err := storage.Store.CountUserLogins(ctx, tx, user.ID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Не смог получить авторизацию пользователя: %v", err))
	}

	dataVersion++
	meta := authlogin.GetMetaInfo(ctx)

	history := &models.SyncHistory{
		UserID:        user.ID,
		SyncTimestamp: time.Now(),
		DataVersion:   dataVersion,
		OperationType: "login",
		MetaInfo:      fmt.Sprintf("Client IP: %s, User-Agent: %s", meta.ClientIP, meta.UserAgent),
	}

	_, err = storage.Store.AddSyncHistory(ctx, tx, history)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Ошибки при добавлении истории: %v", err))
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.New(fmt.Sprintf("Ошибка фиксации транзакции: %v", err))
	}
	return user, nil
}
