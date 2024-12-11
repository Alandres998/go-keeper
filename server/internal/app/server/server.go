package server

import (
	"net"

	"github.com/Alandres998/go-keeper/proto/auth"
	"github.com/Alandres998/go-keeper/proto/options"
	"github.com/Alandres998/go-keeper/proto/private"
	"github.com/Alandres998/go-keeper/server/internal/app/db/storagefactory"
	logger "github.com/Alandres998/go-keeper/server/internal/app/loger"
	v1 "github.com/Alandres998/go-keeper/server/internal/app/routing/v1"
	syncmanager "github.com/Alandres998/go-keeper/server/internal/app/sync"
	configserver "github.com/Alandres998/go-keeper/server/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	syncManager *syncmanager.SyncManager[*private.PrivateDataSyncResponse]
}

func NewServer(syncManager *syncmanager.SyncManager[*private.PrivateDataSyncResponse]) *Server {
	return &Server{syncManager: syncManager}
}

// RunServer Запускаем сервер
func RunServer() {
	configserver.InitConfig()
	storagefactory.NewStorage()

	grpcServer := grpc.NewServer()

	authService := &v1.AuthService{}
	auth.RegisterAuthServiceServer(grpcServer, authService)

	optionsService := &v1.OptionsServiceServer{}
	options.RegisterOptionsServiceServer(grpcServer, optionsService)

	globalSyncManager := syncmanager.NewSyncManager[*private.PrivateDataSyncResponse]()
	privateService := v1.NewPrivateServiceServer(globalSyncManager)
	private.RegisterPrivateServiceServer(grpcServer, privateService)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", configserver.Options.ServerAdress)
	if err != nil {
		logger.LogError("Ошибка создания TCP слушателя", err.Error())
		return
	}

	logger.LoginInfo("gRPC сервер запущен", configserver.Options.ServerAdress)

	if err := grpcServer.Serve(listener); err != nil {
		logger.LogError("Ошибка при запуске gRPC сервера", err.Error())
	}
}
