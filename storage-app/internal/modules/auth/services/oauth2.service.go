package services

import (
	"context"
)

type UserInfo interface {
	GetEmail() string
	GetFirstName() string
	GetLastName() string
	GetImage() string
	GetAuthProvider() string
}

type Oauth2Service interface {
	ExchangeToken(context.Context, string) (UserInfo, error)
}

type Oauth2ServiceToken int

const (
	GoogleOauth2Token Oauth2ServiceToken = iota + 1
	GithubOauth2Token
)

type Oauth2ServiceFactory interface {
	GetOauth2Service(Oauth2ServiceToken) Oauth2Service
	Register(Oauth2ServiceToken, Oauth2Service)
}
type Oauth2ServiceFactoryImpl struct {
	oauth2Services map[Oauth2ServiceToken]Oauth2Service
}

func (o *Oauth2ServiceFactoryImpl) GetOauth2Service(oauth2Token Oauth2ServiceToken) Oauth2Service {
	for token, oauth2Service := range o.oauth2Services {
		if token == oauth2Token {
			return oauth2Service
		}
	}
	return nil
}
func (o *Oauth2ServiceFactoryImpl) Register(token Oauth2ServiceToken, service Oauth2Service) {
	o.oauth2Services[token] = service
}

func NewOauth2ServiceFactory() Oauth2ServiceFactory {
	return &Oauth2ServiceFactoryImpl{
		oauth2Services: make(map[Oauth2ServiceToken]Oauth2Service),
	}
}
