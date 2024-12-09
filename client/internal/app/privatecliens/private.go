package privatecliens

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Alandres998/go-keeper/client/internal/app/sync"
	configclient "github.com/Alandres998/go-keeper/client/internal/config"
	"github.com/Alandres998/go-keeper/proto/private"
	"github.com/manifoldco/promptui"
	"google.golang.org/grpc"
)

// Client Структура клиентаНаСокетах
type Client struct {
	client private.PrivateServiceClient
	conn   *grpc.ClientConn
}

// SendPrivateDataClient Закинуть данные на сервер
func SendPrivateDataClient(conn *grpc.ClientConn) (*private.FillPrivateDataResponse, error) {
	var cardNumber string
	var textData string
	var binaryData []byte

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Номер карты: ")
	if scanner.Scan() {
		cardNumber = scanner.Text()
	}

	fmt.Print("Описание: ")
	if scanner.Scan() {
		textData = scanner.Text()
	}

	fmt.Print("Загружаем файл в бинарник: ")
	if scanner.Scan() {
		binaryData = []byte(scanner.Text())
	}

	ctx, cancel := context.WithTimeout(context.Background(), configclient.Options.TimeOut)
	defer cancel()

	client := private.NewPrivateServiceClient(conn)

	req := &private.FillPrivateDataRequest{
		CardNumber: cardNumber,
		TextData:   textData,
		BinaryData: binaryData,
		Token:      configclient.Options.UserToken,
	}

	resp, err := client.FillPrivateData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вызове метода FillPrivateData: %v", err)
	}
	sync.SyncPrivateData(conn)
	return resp, nil
}

// PrivateDataSync получить даныне из сервера
func PrivateDataSync(conn *grpc.ClientConn) error {
	// Создаём запрос
	ctx, cancel := context.WithTimeout(context.Background(), configclient.Options.TimeOut)
	defer cancel()

	client := private.NewPrivateServiceClient(conn)
	// Открываем поток для общения с сервером
	stream, err := client.SyncPrivateData(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать потоковый запрос: %v", err)
	}

	err = stream.Send(&private.PrivateDataSyncRequest{
		Token: configclient.Options.UserToken,
	})
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	// Получаем ответы от сервера
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			// Поток завершён
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка при получении данных: %v", err)
		}
		configclient.Options.PrivatData.CardNumber = resp.CardNumber
		configclient.Options.PrivatData.TextData = resp.TextData
		configclient.Options.PrivatData.BinaryData = resp.BinaryData
		configclient.Options.PrivatData.MetaInfo = resp.MetaInfo
		configclient.Options.PrivatData.UpdatedAt, err = time.Parse("2006-01-02T15:04:05Z", resp.UpdatedAt)
		if err != nil {
			return fmt.Errorf("не смог отформатировать время: %v", err)
		}
	}

	return nil
}

// LaunchPrivateData функция запуска опросника
func LaunchPrivateData(conn *grpc.ClientConn) {
	exit := false
	// Запускаем горутину для фонового обновления данных
	go sync.SyncPrivateData(conn)

	for {
		PrivateDataSync(conn)
		if configclient.Options.PrivatData.CardNumber == "" {
			SendPrivateDataClient(conn)
		}

		prompt := promptui.Select{
			Label: "Выберите опцию",
			Items: []string{"Изменить данные", "Выйти"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			log.Fatalf("Не правильный текст: %v", err)
		}
		switch result {
		case "Изменить данные":
			SendPrivateDataClient(conn)
		case "Выйти":
			exit = true
		}
		if exit {
			break
		}
	}
}
