package services

import (
	"context"
	"fmt"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GenerateAccessToken(context.Context, string) (string, error)
	GenerateRefreshToken(context.Context, string) (string, error)
	VerifyAccessToken(context.Context, string) (*TokenClaims, error)
	VerifyRefreshToken(context.Context, string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

type JwtServiceImpl struct {
	accessTokenSecret  string
	accessTokenAge     int
	refreshTokenSecret string
	refreshTokenAge    int
	logger             logger.Logger
	isAdmin            bool
}

func verifyToken(token string, secret string) (*TokenClaims, error) {
	tokenDecode, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tokenDecode.Claims.(*TokenClaims); ok && tokenDecode.Valid {
		if claims.ExpiresAt != nil {
			expireTime := claims.ExpiresAt.Time
			if time.Now().After(expireTime) {
				return nil, exception.ErrTokenExpire
			}
			return &TokenClaims{
				UserId: claims.UserId,
			}, nil
		}
	}
	return nil, exception.ErrInvalidToken
}

func generateToken(userId string, secret string, tokenAge int) (string, error) {
	claims := &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenAge) * time.Hour)),
		},
		UserId: userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JwtServiceImpl) VerifyRefreshToken(_ context.Context, token string) (*TokenClaims, error) {
	return verifyToken(token, j.refreshTokenSecret)
}

func (j *JwtServiceImpl) VerifyAccessToken(_ context.Context, token string) (*TokenClaims, error) {
	return verifyToken(token, j.accessTokenSecret)
}

func (j *JwtServiceImpl) GenerateAccessToken(_ context.Context, userId string) (string, error) {
	return generateToken(userId, j.accessTokenSecret, j.accessTokenAge)
}

func (j *JwtServiceImpl) GenerateRefreshToken(_ context.Context, userId string) (string, error) {
	return generateToken(userId, j.refreshTokenSecret, j.refreshTokenAge)
}

func NewUserJwtService(cfg config.JwtConfig, logger logger.Logger) JwtService {
	return &JwtServiceImpl{
		accessTokenSecret:  cfg.AccessToken.Secret,
		accessTokenAge:     cfg.AccessToken.Age,
		refreshTokenSecret: cfg.RefreshToken.Secret,
		refreshTokenAge:    cfg.RefreshToken.Age,
		logger:             logger,
		isAdmin:            false,
	}
}
