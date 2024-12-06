package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/Alandres998/go-keeper/proto/private"
	"github.com/Alandres998/go-keeper/server/internal/app/auth" // Импортируем пакет для валидации токенов
	"github.com/Alandres998/go-keeper/server/internal/app/db/storage"
	logger "github.com/Alandres998/go-keeper/server/internal/app/loger"
	"github.com/Alandres998/go-keeper/server/internal/app/models"
)

type PrivateServiceServer struct {
	private.UnimplementedPrivateServiceServer
}

// FillPrivateData обрабатывает запрос на добавление личных данных.
func (s *PrivateServiceServer) FillPrivateData(ctx context.Context, req *private.FillPrivateDataRequest) (*private.FillPrivateDataResponse, error) {
	// Проверка наличия токена в запросе
	if req.Token == "" {
		return &private.FillPrivateDataResponse{
			Success: false,
			Message: "Токен не был передан",
		}, nil
	}
	logger.LoginInfo("Приватные данные", "Принял запрос есть токен")
	userID, err := auth.ValidateToken(req.Token)
	if err != nil {
		return &private.FillPrivateDataResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка валидации токена: %v", err),
		}, nil
	}
	logger.LoginInfo("Приватные данные", "Токен валиден")
	if req.CardNumber == "" {
		return &private.FillPrivateDataResponse{
			Success: false,
			Message: "Поле card_number не должно быть пустым",
		}, nil
	}
	meta := auth.GetMetaInfo(ctx)

	privateData := &models.PrivateData{
		UserID:     userID,
		DataType:   "default",
		DataKey:    "key",
		TextData:   req.TextData,
		BinaryData: req.BinaryData,
		CardNumber: req.CardNumber,
		MetaInfo:   fmt.Sprintf(`{"source": "grpc", "time": "%s" ,"meta:" "%s"}`, time.Now().Format(time.RFC3339), meta),
	}
	logger.LoginInfo("Приватные данные", "Заполнил данные")
	tx, err := storage.Store.BeginTx(ctx)
	if err != nil {
		return &private.FillPrivateDataResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка начала транзакции: %v", err),
		}, nil
	}

	dataVersion, err := storage.Store.CountUserLogins(ctx, tx, userID)
	if err != nil {
		return &private.FillPrivateDataResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка при синхронизации: %v", err),
		}, nil
	}

	dataVersion++

	history := &models.SyncHistory{
		UserID:        userID,
		SyncTimestamp: time.Now(),
		DataVersion:   dataVersion,
		OperationType: "login",
		MetaInfo:      fmt.Sprintf("Client IP: %s, User-Agent: %s", meta.ClientIP, meta.UserAgent),
	}
	_, err = storage.Store.AddSyncHistory(ctx, tx, history)
	logger.LoginInfo("Приватные данные", "Синхронизировал данные")
	data, err := storage.Store.InsertPrivateData(ctx, tx, privateData)
	if err != nil {
		tx.Rollback()
		return &private.FillPrivateDataResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка сохранения данных: %v", err),
		}, nil
	}

	err = tx.Commit()
	logger.LoginInfo("Приватные данные", "Коммит")
	if err != nil {
		return &private.FillPrivateDataResponse{
			Success: false,
			Message: fmt.Sprintf("Ошибка завершения транзакции: %v", err),
		}, nil
	}

	return &private.FillPrivateDataResponse{
		Success: true,
		Message: fmt.Sprintf("Данные успешно сохранены. ID: %d", data.ID),
	}, nil
}

// GetPrivateData возвращает данные из таблицы
func (s *PrivateServiceServer) GetPrivateData(ctx context.Context, req *private.GetPrivateDataRequest) (*private.GetPrivateDataResponse, error) {
	// Проверяем токен
	userID, err := auth.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("Недействительный токен")
	}

	// Получаем данные из таблицы
	data, err := storage.Store.GetPrivateDataByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения данных: %v", err)
	}

	// Формируем ответ
	var privateDataList []*private.PrivateData
	for _, d := range data {
		privateDataList = append(privateDataList, &private.PrivateData{
			Id:         int32(d.ID),
			CardNumber: d.CardNumber,
			TextData:   d.TextData,
			BinaryData: d.BinaryData,
			MetaInfo:   d.MetaInfo,
		})
	}

	return &private.GetPrivateDataResponse{
		PrivateData: privateDataList,
	}, nil
}
