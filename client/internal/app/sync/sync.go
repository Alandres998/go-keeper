package sync

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/Alandres998/go-keeper/client/internal/app/consoleclient"
	userservices "github.com/Alandres998/go-keeper/client/internal/app/services/userServices"
	configclient "github.com/Alandres998/go-keeper/client/internal/config"
	"github.com/Alandres998/go-keeper/proto/private"
	"google.golang.org/grpc"
)

// SyncPrivateData функция запуска опросника
func SyncPrivateData(conn *grpc.ClientConn) {
	client := private.NewPrivateServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаём поток для получения обновлений
	stream, err := client.SyncPrivateData(ctx)
	if err != nil {
		log.Printf("Ошибка при создании потока: %v", err)
		return
	}

	// Отправляем токен пользователя для авторизации
	err = stream.Send(&private.PrivateDataSyncRequest{
		Token: configclient.Options.UserToken,
	})
	if err != nil {
		log.Printf("Ошибка при отправке данных в поток: %v", err)
		return
	}

	log.Println("Подписка на обновления успешна")

	// Получаем данные из потока
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("Поток завершён")
			break
		}
		if err != nil {
			log.Printf("Ошибка при получении данных: %v", err)
			break
		}

		configclient.Options.PrivatData.CardNumber = resp.CardNumber
		configclient.Options.PrivatData.TextData = resp.TextData
		configclient.Options.PrivatData.BinaryData = resp.BinaryData
		configclient.Options.PrivatData.MetaInfo = resp.MetaInfo
		updatedAt, err := time.Parse(time.RFC3339, resp.UpdatedAt)
		if err != nil {
			log.Printf("Не смог отформатировать время: %v", err)
			continue
		}

		configclient.Options.PrivatData.UpdatedAt = updatedAt

		consoleclient.ClearConsole()
		// Отобразить обновлённые данные
		userservices.PrintPrivateInfo()
	}
}
