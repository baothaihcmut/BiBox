package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

type FacebookUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture struct {
		Data struct {
			Height       int    `json:"height"`
			IsSilhouette bool   `json:"is_silhouette"`
			URL          string `json:"url"`
			Width        int    `json:"width"`
		} `json:"data"`
	} `json:"picture"`
}

func (g *FacebookUser) GetEmail() string {
	return g.Email
}
func (g *FacebookUser) GetFirstName() string {
	return strings.Split(g.Name, " ")[0]
}
func (g *FacebookUser) GetLastName() string {
	return strings.Split(g.Name, " ")[1]
}
func (g *FacebookUser) GetImage() string {
	return g.Picture.Data.URL
}

type FacebookOauth2Service struct {
	facbookOauth2Config *oauth2.Config
}

func (o *FacebookOauth2Service) ExchangeToken(ctx context.Context, authCode string) (UserInfo, error) {
	token, err := o.facbookOauth2Config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(fmt.Sprintf("https://graph.facebook.com/v18.0/me?fields=id,name,email,picture.width(300).height(300)&access_token=%s", token.AccessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user FacebookUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func NewFacebookOauth2Service(oauth2 *oauth2.Config) Oauth2Service {
	return &FacebookOauth2Service{
		facbookOauth2Config: oauth2,
	}
}
