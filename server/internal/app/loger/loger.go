package logger

import (
	"log"

	"go.uber.org/zap"
)

// LoginInfo записывает информационное сообщение в лог.
// title - заголовок сообщения, info - сообщение
func LoginInfo(title string, info string) {
	// Создание нового логгера
	logger, errLog := zap.NewProduction()
	if errLog != nil {
		log.Fatalf("Не смог инициализировать логгер: %v", errLog)
	}

	// Используем defer для синхронизации логгера
	defer func() {
		_ = logger.Sync()
	}()

	// Записываем информационное сообщение в лог
	logger.Info("Внимание",
		zap.String(title, info),
	)
}

// LogError записывает сообщение об ошибке в лог.
// title - заголовок сообщения об ошибке, info - сообщение
func LogError(title string, info string) {
	logger, errLog := zap.NewProduction()

	defer func() {
		_ = logger.Sync()
	}()

	if errLog != nil {
		log.Fatalf("Не смог инициализировать логгер")
	}

	logger.Error("Ошибка",
		zap.String(title, info),
	)
}
