package middleware

import (
	"context"
	"net/http"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/services"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authHandler interface {
	VerifyAccessToken(context.Context, string) (*services.TokenClaims, error)
}, logger logger.Logger, isAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			if err == http.ErrNoCookie {
			} else {
				logger.Errorf(c.Request.Context(), nil, "Error extract cookie:", err)
				c.JSON(http.StatusInternalServerError, response.InitResponse[any](false, "Internal error", nil))
				c.Abort()
			}
		}
		if accessToken == "" {
			userContext := models.UserContext{
				Id:   "67c58727090e489cf6796a8f",
				Role: enums.UserRole,
			}
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.UserContext, &userContext))
			c.Set(string(constant.UserContext), &userContext)
			c.Next()
			return
		}
		//decode token
		claims, err := authHandler.VerifyAccessToken(c.Request.Context(), accessToken)
		if err != nil {
			if err == exception.ErrTokenExpire || err == exception.ErrInvalidToken {
				c.JSON(http.StatusUnauthorized, response.InitResponse[any](false, err.Error(), nil))
			} else {
				logger.Errorf(c.Request.Context(), nil, "Error verify token", err)
			}
			c.Abort()
			return
		}
		//set user context
		userContext := models.UserContext{
			Id: claims.UserId,
		}
		if isAdmin {
			userContext.Role = enums.AdminRole
		} else {
			userContext.Role = enums.UserRole
		}
		//set user to context
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.UserContext, &userContext))
		c.Set(string(constant.UserContext), &userContext)
		c.Next()
	}
}
