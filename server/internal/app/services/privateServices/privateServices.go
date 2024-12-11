package privateservices

import (
	"context"
	"fmt"

	"github.com/Alandres998/go-keeper/models"
	"github.com/Alandres998/go-keeper/server/internal/app/auth"
	"github.com/Alandres998/go-keeper/server/internal/app/db/storage"
	configserver "github.com/Alandres998/go-keeper/server/internal/config"
)

func CreatePrivateData(ctx context.Context, userID int, textData string, binaryData []byte, cardNumber string) (*models.PrivateData, error) {
	meta := auth.GetMetaInfo(ctx)
	// Сохраняем данные
	privateData := &models.PrivateData{
		UserID:     userID,
		TextData:   textData,
		BinaryData: binaryData,
		CardNumber: cardNumber,
		MetaInfo:   fmt.Sprintf("ClientIP: %s, UserAgent: %s", meta.ClientIP, meta.UserAgent),
	}

	ctx, cancel := context.WithTimeout(ctx, configserver.Options.TimeOut)
	defer cancel()

	tx, err := storage.Store.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err = storage.Store.InsertOrUpdatePrivateData(ctx, tx, privateData); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return privateData, nil
}
