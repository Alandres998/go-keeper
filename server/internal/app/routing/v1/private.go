package v1

import (
	"context"
	"strconv"
	"time"

	"github.com/Alandres998/go-keeper/models"
	"github.com/Alandres998/go-keeper/proto/private"
	"github.com/Alandres998/go-keeper/server/internal/app/auth" // Импортируем пакет для валидации токенов
	"github.com/Alandres998/go-keeper/server/internal/app/db/storage"
	syncmanager "github.com/Alandres998/go-keeper/server/internal/app/sync"
)

// Структура для хранения соединений
type PrivateServiceServer struct {
	private.UnimplementedPrivateServiceServer
	syncManager *syncmanager.SyncManager[*private.PrivateDataSyncResponse]
}

// Иницилазиация сервера
func NewPrivateServiceServer(syncManager *syncmanager.SyncManager[*private.PrivateDataSyncResponse]) *PrivateServiceServer {
	return &PrivateServiceServer{
		syncManager: syncManager,
	}
}

func (s *PrivateServiceServer) SyncPrivateData(stream private.PrivateService_SyncPrivateDataServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	userID, err := auth.ValidateToken(req.GetToken())
	if err != nil {
		return err
	}

	s.syncManager.AddStream(strconv.Itoa(userID), stream)
	defer s.syncManager.RemoveStream(strconv.Itoa(userID), stream)

	data, err := storage.Store.GetPrivateDataByUserID(stream.Context(), userID)
	if err != nil {
		return err
	}
	if data != nil {
		resp := &private.PrivateDataSyncResponse{
			CardNumber: data.CardNumber,
			TextData:   data.TextData,
			BinaryData: data.BinaryData,
			MetaInfo:   data.MetaInfo,
			UpdatedAt:  time.Now().Format(time.RFC3339),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}

	// Поддерживаем соединение до завершения
	<-stream.Context().Done()
	return nil
}

func (s *PrivateServiceServer) FillPrivateData(ctx context.Context, req *private.FillPrivateDataRequest) (*private.FillPrivateDataResponse, error) {
	if req.Token == "" {
		return &private.FillPrivateDataResponse{Success: false, Message: "Токен отсутствует"}, nil
	}

	// Валидируем токен
	userID, err := auth.ValidateToken(req.Token)
	if err != nil {
		return &private.FillPrivateDataResponse{Success: false, Message: "Ошибка валидации токена"}, nil
	}

	// Сохраняем данные
	privateData := &models.PrivateData{
		UserID:     userID,
		TextData:   req.TextData,
		BinaryData: req.BinaryData,
		CardNumber: req.CardNumber,
		MetaInfo:   "", // Пример, заполняйте метаинформацию, если нужно
	}
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

	// Формируем обновление
	resp := &private.PrivateDataSyncResponse{
		CardNumber: req.CardNumber,
		TextData:   req.TextData,
		BinaryData: req.BinaryData,
		MetaInfo:   "", // Пример, можно добавить актуальную метаинформацию
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}

	// Отправляем обновление всем потокам пользователя
	s.syncManager.Broadcast(strconv.Itoa(userID), resp)

	return &private.FillPrivateDataResponse{Success: true, Message: "Данные успешно сохранены"}, nil
}
