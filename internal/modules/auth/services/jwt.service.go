package services

import (
	"context"
	"fmt"
	"time"

	exception "github.com/baothaihcmut/Storage-app/internal/common/exeption"
	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/golang-jwt/jwt/v5"
)

type TokenSub struct {
	UserId string
}

type JwtService interface {
	GenerateAccessToken(context.Context, string) (string, error)
	GenerateRefreshToken(context.Context, string) (string, error)
	VerifyAccessToken(context.Context, string) (*TokenSub, error)
	VerifyRefreshToken(context.Context, string) (*TokenSub, error)
}

type TokenClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

type JwtServiceImpl struct {
	accessTokenSecret  string
	accessTokenAge     string
	refreshTokenSecret string
	refreshTokenAge    string
	logger             logger.Logger
}

func (j *JwtServiceImpl) VerifyRefreshToken(ctx context.Context, token string) (*TokenSub, error) {

	tokenDecode, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.refreshTokenSecret), nil
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
			return &TokenSub{
				UserId: claims.UserId,
			}, nil
		}
	}
	return nil, exception.ErrInvalidToken
}
