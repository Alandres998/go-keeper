package configclient

import (
	"flag"
	"os"
	"time"

	"github.com/Alandres998/go-keeper/models"
)

// OptionsStruct структура с настройками
type OptionsStruct struct {
	ServerAdress string
	UserToken    string
	TimeOut      time.Duration
	TimeRefresh  time.Duration
	PrivatData   models.PrivateData
}

// Options общая конфигурация проекта
var Options OptionsStruct

// InitConfig инициализация конфига
func InitConfig() {
	Options.TimeOut = 5 * time.Second
	Options.TimeRefresh = 2 * time.Second
	loadEnv()
}

// parseFlags Устанавливаем конфиг из флагов командой строки
func parseFlags() {
	flag.StringVar(&Options.ServerAdress, "h", "http://localhost:8000", "server adress")
	flag.Parse()
}

// loadEnv Устанавливаем конфиг из env
func loadEnv() {
	setOptionIfEmpty(&Options.ServerAdress, os.Getenv("SERVER_ADDRESS"))
}

// setOptionIfEmpty Устанавливает значение, если оно пустое для string
func setOptionIfEmpty(target *string, value string) {
	if *target == "" {
		*target = value
	}
}

// setOptionIfEmptyBool Устанавливает значение, если оно пустое для bool
func setOptionIfEmptyBool(target *bool, value bool) {
	if !*target {
		*target = value
	}
}

// setOptionIfEmptyInt Устанавливает значение, если оно пустое для int
func setOptionIfEmptyInt(target *int, value int) {
	if *target == 0 {
		*target = value
	}
}

// stringToBool Костыль для env с целью преобразования текста в bool
func stringToBool(value string) bool {
	return value == "true"
}

// getEnv Енв с дефолтным значением
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
