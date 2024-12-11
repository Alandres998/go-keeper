package v1

import (
	"context"
	"fmt"

	"github.com/Alandres998/go-keeper/proto/auth"
	authlogin "github.com/Alandres998/go-keeper/server/internal/app/auth"
	"github.com/Alandres998/go-keeper/server/internal/app/services/userServices"
)

type AuthService struct {
	auth.UnimplementedAuthServiceServer // встраиваем не реализованные методы
}

// Register GRPC регистрация
func (s *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	createdUser, err := userServices.UserCreate(ctx, req.GetLogin(), req.GetPassword(), req.GetEmail())

	if err != nil {
		return &auth.RegisterResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка создания пользователя: %v", err),
		}, err
	}

	// Генерация токена
	token, err := authlogin.GenerateToken(ctx, createdUser.ID)
	if err != nil {
		return &auth.RegisterResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка генерации токена: %v", err),
		}, err
	}

	return &auth.RegisterResponse{
		Success: true,
		Message: "Вы успешно зарегистрировались",
		Token:   token,
	}, nil
}

// Login GRPC авторизация
func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := userServices.UserLogin(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}
	token, err := authlogin.GenerateToken(ctx, user.ID)

	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: fmt.Sprintf("Не корректный токен: %v", err),
		}, err
	}

	return &auth.LoginResponse{
		Success: true,
		Message: "Успешная авторизация",
		Token:   token,
	}, nil
}
