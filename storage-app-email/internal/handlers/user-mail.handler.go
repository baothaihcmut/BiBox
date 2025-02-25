package handlers

import "github.com/baothaihcmut/BiBox/storage-app-email/internal/router"

type UserHandler interface {
	Init(r *router.MessageRouter)
}

type UserHandlerImpl struct {
}
