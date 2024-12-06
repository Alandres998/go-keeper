package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Alandres998/go-keeper/client/internal/app/authclient"
	"github.com/Alandres998/go-keeper/client/internal/app/initclient"
	"github.com/Alandres998/go-keeper/client/internal/app/privatecliens"
	configclient "github.com/Alandres998/go-keeper/client/internal/config"
	"google.golang.org/grpc"
)

var (
	buildVersion = "1.0.0"
	buildDate    = "2024-12-04"
	buildCommit  = "abcdef1234567890"
)

func main() {
	configclient.InitConfig()
	initclient.GetInfoAboutVersion()

	ctx, cancel := context.WithTimeout(context.Background(), configclient.Options.TimeOut)
	defer cancel()

	// Подключаемся к серверу
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("localhost%s", configclient.Options.ServerAdress), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Не смог подключиться к серверу: %v", err)
	}
	defer conn.Close()

	err = initclient.CheckConnect(conn)
	if err != nil {
		log.Fatalf("Не смог подключиться к серверу: %v", err)
	}
	fmt.Print("Успешная установка соединения с сервером\n")
	authclient.StartSession(conn)
	privatecliens.LaunchPrivateData(conn)
	//privatecliens.FillPrivateDataClient(conn)
}
