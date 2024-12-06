package privatecliens

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Alandres998/go-keeper/client/internal/app/consoleclient"
	configclient "github.com/Alandres998/go-keeper/client/internal/config"
	"github.com/Alandres998/go-keeper/proto/private"
	"github.com/manifoldco/promptui"
	"google.golang.org/grpc"
)

// Структура клиентаНаСокетах
type ClientSoket struct {
	client private.PrivateServiceClient
	conn   *grpc.ClientConn
}

// Закинуть данные на сервер
func FillPrivateDataClient(conn *grpc.ClientConn) (*private.FillPrivateDataResponse, error) {
	var cardNumber string
	var textData string
	var binaryData []byte

	fmt.Print("Номер карты: ")
	fmt.Scanln(&cardNumber)
	fmt.Print("Описание: ")
	fmt.Scanln(&textData)
	fmt.Print("Загружаем файл в бинарник: ")
	fmt.Scanln(&binaryData)

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
	go syncPrivateDataPeriodically(conn)
	return resp, nil
}

// PrivateDataSync вызывает серверный метод PrivateDataSync
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
	go syncPrivateDataPeriodically(conn)

	for {
		PrivateDataSync(conn)
		if configclient.Options.PrivatData.CardNumber == "" {
			FillPrivateDataClient(conn)
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
			FillPrivateDataClient(conn)
		case "Выйти":
			exit = true
		}
		if exit {
			break
		}
	}
}

// LaunchPrivateData функция запуска опросника
func syncPrivateDataPeriodically(conn *grpc.ClientConn) {
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

		// Обновляем локальные данные
		configclient.Options.PrivatData.CardNumber = resp.CardNumber
		configclient.Options.PrivatData.TextData = resp.TextData
		configclient.Options.PrivatData.BinaryData = resp.BinaryData
		configclient.Options.PrivatData.MetaInfo = resp.MetaInfo

		// Парсинг времени
		updatedAt, err := time.Parse(time.RFC3339, resp.UpdatedAt)
		if err != nil {
			log.Printf("Не смог отформатировать время: %v", err)
			continue
		}
		configclient.Options.PrivatData.UpdatedAt = updatedAt
		consoleclient.ClearConsole()
		// Отобразить обновлённые данные
		PrintPrivateInfo()
	}
}

// Отобразить пользователю свою информацию
func PrintPrivateInfo() {
	fmt.Print("----------------------------------------------\n")
	fmt.Printf("Card Number: %s\n", configclient.Options.PrivatData.CardNumber)
	fmt.Printf("Text Data: %s\n", configclient.Options.PrivatData.TextData)
	fmt.Printf("Binary Data: %v\n", configclient.Options.PrivatData.BinaryData)
	fmt.Printf("Meta Info: %s\n", configclient.Options.PrivatData.MetaInfo)
	fmt.Printf("Updated At: %s\n", configclient.Options.PrivatData.UpdatedAt)
	fmt.Print("----------------------------------------------\n\n\n\n")
}
