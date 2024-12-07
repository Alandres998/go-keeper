package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/Alandres998/go-keeper/models"
	"github.com/Alandres998/go-keeper/proto/auth"
	authlogin "github.com/Alandres998/go-keeper/server/internal/app/auth"
	"github.com/Alandres998/go-keeper/server/internal/app/db/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	auth.UnimplementedAuthServiceServer // встраиваем не реализованные методы
}

func (s *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	// Начинаем транзакцию
	tx, err := storage.Store.BeginTx(ctx)
	if err != nil {
		return &auth.RegisterResponse{
			Success: false,
			Message: "Ошибка при создании транзакции",
		}, err
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
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	}

	// Используем транзакцию при создании пользователя
	createdUser, err := storage.Store.CreateUser(ctx, tx, &user)
	if err != nil {
		return &auth.RegisterResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка при регистрации пользователя: %v", err),
		}, err
	}

	// Генерация токена
	token, err := authlogin.GenerateToken(createdUser.ID)
	if err != nil {
		return &auth.RegisterResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка генерации токена: %v", err),
		}, err
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
		return &auth.RegisterResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибки при добавлении истории: %v", err),
		}, err
	}

	// Коммит транзакции
	if err = tx.Commit(); err != nil {
		return &auth.RegisterResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка фиксации транзакции: %v", err),
		}, err
	}

	return &auth.RegisterResponse{
		Success: true,
		Message: "Вы успешно зарегистрировались",
		Token:   token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	tx, err := storage.Store.BeginTx(ctx)
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: "Ошибка при создании транзакции",
		}, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	user, err := storage.Store.GetUserByLogin(ctx, tx, req.GetLogin())
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: "Не корректный логин или пароль",
		}, nil
	}

	// Сравниваем хешированный пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword()))
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: "Не корректный логин или пароль",
		}, nil
	}

	token, err := authlogin.GenerateToken(user.ID)
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: fmt.Sprintf("Не корректный токен: %v", err),
		}, err
	}

	dataVersion, err := storage.Store.CountUserLogins(ctx, tx, user.ID)
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: fmt.Sprintf("Не смог получить авторизацию пользователя: %v", err),
		}, err
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
		return &auth.LoginResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибки при добавлении истории: %v", err),
		}, err
	}

	if err = tx.Commit(); err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка фиксации транзакции: %v", err),
		}, err
	}

	return &auth.LoginResponse{
		Success: true,
		Message: "Успешная авторизация",
		Token:   token,
	}, nil
}
