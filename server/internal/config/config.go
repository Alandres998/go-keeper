package configserver

import (
	"flag"
	"os"
)

// OptionsStruct структура с настройками
type OptionsStruct struct {
	ServerAdress string
	DatabaseDSN  string
	SecretJWTKey string
}

// Options общая конфигурация проекта
var Options OptionsStruct

// InitConfig инициализация конфига
func InitConfig() {
	loadEnv()
}

// parseFlags Устанавливаем конфиг из флагов командой строки
func parseFlags() {
	flag.StringVar(&Options.ServerAdress, "h", "http://localhost:8080", "server adress")
	flag.StringVar(&Options.DatabaseDSN, "d", "http://localhost:8080", "data base dns")
	flag.StringVar(&Options.ServerAdress, "k", "secret key", "secret key")
	flag.Parse()
}

// loadEnv Устанавливаем конфиг из env
func loadEnv() {
	setOptionIfEmpty(&Options.ServerAdress, os.Getenv("SERVER_ADDRESS"))
	setOptionIfEmpty(&Options.DatabaseDSN, os.Getenv("DATABASE_DSN"))
	setOptionIfEmpty(&Options.DatabaseDSN, getEnv("SecretJWTKey", "secret"))
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

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
