package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Константы для JWT
const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// Заглушка
var refreshTokenSecret = []byte("1231241012030ggg")

// AccessClaims пользовательские JWT claims
type AccessClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
type RefreshClaims struct {
	TokenID string `json:"token_id"`
	UserID  string `json:"user_id"`
	jwt.StandardClaims
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// GenerateToken создает новый Access JWT токен для пользователя
func GenerateAccessToken(userID string, secretKey string) (string, error) {
	claims := &AccessClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func GenerateRefreshToken(userID, secretKey string) (string, string, error) {
	tokenID := uuid.New().String()
	claims := &RefreshClaims{
		TokenID: tokenID,
		UserID:  userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefreshTokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, tokenID, nil
}

// ValidateToken проверяет валидность JWT токена
func ValidateAccessToken(tokenString string, secretKey string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ValidRefreshToken(ctx context.Context, tokenString, secretKey string, dbResponce func(context.Context, string) bool) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		// Проверяем, существует ли токен в хранилище
		if exists := dbResponce(ctx, claims.TokenID); !exists {
			return nil, fmt.Errorf("token revoked or doesn't exist")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
