package initclient

import (
	"context"
	"fmt"

	"github.com/Alandres998/go-keeper/proto/options"
	"google.golang.org/grpc"
)

var (
	buildVersion = "1.0.0"
	buildDate    = "2024-12-04"
	buildCommit  = "abcdef1234567890"
)

// GetInfoAboutVersion Выводим данные о билде
func GetInfoAboutVersion() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

// Провряем соединение
func CheckConnect(conn *grpc.ClientConn) error {
	client := options.NewOptionsServiceClient(conn)

	// Проверка связи с сервером
	_, err := client.Ping(context.Background(), &options.PingRequest{})
	return err
}
