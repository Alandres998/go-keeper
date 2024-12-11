package v1

import (
	"context"

	"github.com/Alandres998/go-keeper/proto/options"
)

type OptionsServiceServer struct {
	options.UnimplementedOptionsServiceServer
}

// Метод для проверки клиентом доступности сервера
func (s *OptionsServiceServer) Ping(ctx context.Context, req *options.PingRequest) (*options.PingResponse, error) {
	return &options.PingResponse{Message: "Pong"}, nil
}
