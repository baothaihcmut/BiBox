package middleware

import (
	"fmt"
	"net/http"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			fmt.Println(c.Errors)
			status := exception.ErrorStatusMapper(c.Errors[0])
			if status != http.StatusInternalServerError {
				c.JSON(status, response.InitResponse[any](false, c.Errors[0].Error(), nil))
				return
			}
			c.JSON(http.StatusInternalServerError, response.InitResponse[any](false, "Internal error", nil))
		}
	}
}
