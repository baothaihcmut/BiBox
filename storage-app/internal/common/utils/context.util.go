package utils

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
)

func GetUserContext(ctx context.Context) *models.UserContext {
	userContext := ctx.Value(constant.UserContext).(*models.UserContext)
	return userContext

}
