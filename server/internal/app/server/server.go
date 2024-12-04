package server

import (
	"net"

	"github.com/Alandres998/go-keeper/proto/auth"
	"github.com/Alandres998/go-keeper/server/internal/app/db/storagefactory"
	logger "github.com/Alandres998/go-keeper/server/internal/app/loger"
	v1 "github.com/Alandres998/go-keeper/server/internal/app/routing/v1"
	configserver "github.com/Alandres998/go-keeper/server/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer() {
	configserver.InitConfig()
	storagefactory.NewStorage()

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервис аутентификации
	authService := &v1.AuthService{}
	auth.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)

	// Создаем TCP слушатель
	listener, err := net.Listen("tcp", configserver.Options.ServerAdress)
	if err != nil {
		logger.LogError("Ошибка создания TCP слушателя", err.Error())
		return
	}

	// Логирование успешного запуска
	logger.LoginInfo("gRPC сервер запущен", configserver.Options.ServerAdress)

	// Запуск gRPC сервера
	if err := grpcServer.Serve(listener); err != nil {
		logger.LogError("Ошибка при запуске gRPC сервера", err.Error())
	}
}
