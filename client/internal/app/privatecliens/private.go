package privatecliens

import (
	"context"
	"fmt"

	configclient "github.com/Alandres998/go-keeper/client/internal/config"
	"github.com/Alandres998/go-keeper/proto/private"
	"google.golang.org/grpc"
)

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

	return resp, nil
}
