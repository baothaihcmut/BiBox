package middleware

import (
	"context"
	"net/http"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		//generate request id for logging
		requestId := uuid.New()
		c.Set(string(constant.RequestIdContext), requestId.String())
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.RequestIdContext, requestId.String()))
		logger.Debug(c.Request.Context(), map[string]interface{}{
			"uri":    c.Request.URL.String(),
			"method": c.Request.Method,
		}, "Incoming request")
		c.Next()
		status := c.Writer.Status()
		logger.Debug(c.Request.Context(), map[string]interface{}{
			"status": status,
		}, "Outgoing response")
		if status == http.StatusInternalServerError {
			//get user context
			userContext, _ := c.Get(string(constant.UserContext))
			logger.Error(c.Request.Context(), map[string]interface{}{
				"user_id": userContext.(*models.UserContext).Id,
				"detail":  c.Errors,
			}, c.Errors[0].Error())
		}
	}
}
