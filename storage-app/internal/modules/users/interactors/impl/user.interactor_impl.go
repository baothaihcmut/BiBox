package impl

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/repositories"
)

type UserInteractorImpl struct {
	userRepo repositories.UserRepository
}

func NewUserInteractor(
	userRepo repositories.UserRepository,
) interactors.UserInteractor {
	return &UserInteractorImpl{
		userRepo: userRepo,
	}
}
