package middleware

import (
	"net/http"

	"github.com/baothaihcmut/Storage-app/internal/common/constant"
	"github.com/baothaihcmut/Storage-app/internal/common/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateMiddleware[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto T
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, response.InitResponse(false, err.Error(), nil))
			c.Abort()
			return
		}
		if err := validate.Struct(&dto); err != nil {
			c.JSON(http.StatusBadRequest, response.InitResponse(false, err.Error(), nil))
			c.Abort()
			return
		}
		c.Set(string(constant.PayloadContext), &dto)
		c.Next()
	}
}
