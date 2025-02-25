package services

import (
	"context"
	"encoding/json"

	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"golang.org/x/oauth2"
)

type GoogleUserInfo struct {
	Email     string `json:"email"`
	FirstName string `json:"family_name"`
	LastName  string `json:"given_name"`
	Image     string `json:"picture"`
}

func (g *GoogleUserInfo) GetEmail() string {
	return g.Email
}
func (g *GoogleUserInfo) GetFirstName() string {
	return g.FirstName
}
func (g *GoogleUserInfo) GetLastName() string {
	return g.LastName
}
func (g *GoogleUserInfo) GetImage() string {
	return g.Image
}

type GoogleOauth2Service struct {
	oauth2Config *oauth2.Config
	logger       logger.Logger
}

func (o *GoogleOauth2Service) ExchangeToken(ctx context.Context, authCode string) (UserInfo, error) {
	token, err := o.oauth2Config.Exchange(ctx, authCode)
	if err != nil {
		o.logger.Errorf(ctx, map[string]interface{}{
			"authCode": authCode,
		}, "Error exchange google token:", err)
		return nil, err
	}
	o.logger.Debug(ctx, map[string]interface{}{
		"token": token,
	}, "Exchange token success")
	client := o.oauth2Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		o.logger.Errorf(ctx, nil, "Error get user info:", err)
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		o.logger.Errorf(ctx, nil, "Error decode user info:", err)
	}
	o.logger.Debug(ctx, map[string]interface{}{
		"email": userInfo.Email,
	}, "Get user info success")

	return &userInfo, nil
}
func NewGoogleOauth2Service(oauth2 *oauth2.Config, logger logger.Logger) Oauth2Service {
	return &GoogleOauth2Service{
		oauth2Config: oauth2,
		logger:       logger,
	}
}
