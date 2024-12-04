package auth

import (
	"context"
	"fmt"
	"time"

	configserver "github.com/Alandres998/go-keeper/server/internal/config"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Секретный ключ
var secretKey = []byte(configserver.Options.SecretJWTKey)

// Структура для меты
type MetaRequestInfo struct {
	ClientIP  string
	UserAgent string
}

// GenerateToken Сгенерировать токен
func GenerateToken(userID int) (string, error) {
	// Сюда бы конечно можно закинуть payLoad но это же учебный проект
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("не смог сгенерировать ключ: %v", err)
	}

	return tokenString, nil
}

// GetMetaInfo Получить мета информацию
func GetMetaInfo(ctx context.Context) MetaRequestInfo {
	var meta MetaRequestInfo

	meta.ClientIP = "unknown"
	meta.UserAgent = "unknown"

	// Попытка получить IP-адрес клиента
	if peer, ok := peer.FromContext(ctx); ok {
		meta.ClientIP = peer.Addr.String()
	}

	// Попытка извлечь User-Agent из контекста (если передается)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			meta.UserAgent = ua[0]
		}
	}
	return meta
}