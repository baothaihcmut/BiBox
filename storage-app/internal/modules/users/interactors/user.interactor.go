package interactors

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/presenters"
)

type UserInteractor interface {
	SearchUserByEmail(context.Context, *presenters.SearchUserByEmailInput) (*presenters.SearchUserByEmailOuput, error)
}
