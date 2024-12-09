package userservices

import (
	"fmt"

	configclient "github.com/Alandres998/go-keeper/client/internal/config"
)

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
