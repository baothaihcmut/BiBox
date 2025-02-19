package middleware

import (
	"net/http"

	"github.com/baothaihcmut/Storage-app/internal/common/constant"
	"github.com/baothaihcmut/Storage-app/internal/common/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ValidateMiddleware[T any](bindUri bool, bindings ...binding.Binding) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto T
		if bindUri {
			if err := c.ShouldBindUri(&dto); err != nil {
				c.JSON(http.StatusBadRequest, response.InitResponse[any](false, err.Error(), nil))
				c.Abort()
				return
			}
		}
		for _, binding := range bindings {
			if err := c.MustBindWith(&dto, binding); err != nil {
				c.JSON(http.StatusBadRequest, response.InitResponse[any](false, err.Error(), nil))
				c.Abort()
				return
			}
		}

		c.Set(string(constant.PayloadContext), &dto)
		c.Next()
	}
}
