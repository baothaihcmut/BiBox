package services

import (
	"context"
	"encoding/json"

	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	Email     string `json:"email"`
	FirstName string `json:"family_name"`
	LastName  string `json:"given_name"`
	Image     string `json:"picture"`
}

type Oauth2Service interface {
	ExchangeToken(context.Context, string) (*UserInfo, error)
}
type GoogleOauth2Service struct {
	oauth2Config *oauth2.Config
	logger       logger.Logger
}

func (o *GoogleOauth2Service) ExchangeToken(ctx context.Context, authCode string) (*UserInfo, error) {
	token, err := o.oauth2Config.Exchange(ctx, authCode)
	if err != nil {
		o.logger.Errorf(ctx, map[string]interface{}{
			"authCode": authCode,
		}, "Error exchange google token:", err)
		return nil, err
	}
	client := o.oauth2Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		o.logger.Errorf(ctx, nil, "Error get user info:", err)
	}
	defer resp.Body.Close()
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		o.logger.Errorf(ctx, nil, "Error decode user info:", err)
	}
	return &userInfo, nil
}
