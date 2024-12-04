package storagefactory

import (
	"log"

	"github.com/Alandres998/go-keeper/server/internal/app/db/db"
	"github.com/Alandres998/go-keeper/server/internal/app/db/storage"
	configserver "github.com/Alandres998/go-keeper/server/internal/config"
	"go.uber.org/zap"
)

// NewStorage фабрика на всякий если будут сюрпризы со сторами
func NewStorage() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	storage.Store, err = db.NewDBStorage(configserver.Options.DatabaseDSN)
	if err != nil {
		logger.Error("Не удалось иницировать хранилище",
			zap.String("Ошибка", err.Error()),
		)
		return
	}
}
