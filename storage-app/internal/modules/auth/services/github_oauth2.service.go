package services

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"golang.org/x/oauth2"
)

type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type GithubOauth2Service struct {
	oauth2Config *oauth2.Config
	logger       logger.Logger
}

func (g *GitHubUserInfo) GetFirstName() string {
	parts := strings.Fields(g.Name) // Split by spaces
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func (g *GitHubUserInfo) GetLastName() string {
	parts := strings.Fields(g.Name)
	if len(parts) > 1 {
		return strings.Join(parts[1:], " ") // Join remaining parts as last name
	}
	return ""
}

func (g *GitHubUserInfo) GetEmail() string {
	return g.Email
}

func (g *GitHubUserInfo) GetImage() string {
	return g.AvatarURL
}

func (g *GitHubUserInfo) GetAuthProvider() string {
	return "github"
}

func (o *GithubOauth2Service) ExchangeToken(ctx context.Context, authCode string) (UserInfo, error) {
	token, err := o.oauth2Config.Exchange(ctx, authCode)
	if err != nil {
		o.logger.Errorf(ctx, map[string]interface{}{
			"authCode": authCode,
		}, "Error exchanging GitHub token:", err)
		return nil, err
	}

	o.logger.Debug(ctx, map[string]interface{}{
		"token": token.AccessToken,
	}, "Exchange token success")

	client := o.oauth2Config.Client(ctx, token)

	// Get user info from GitHub API
	userInfoURL := "https://api.github.com/user"
	resp, err := client.Get(userInfoURL)
	if err != nil {
		o.logger.Errorf(ctx, nil, "Error getting GitHub user info:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		o.logger.Errorf(ctx, nil, "Error reading GitHub user info:", err)
		return nil, err
	}

	var user GitHubUserInfo
	if err := json.Unmarshal(body, &user); err != nil {
		o.logger.Errorf(ctx, nil, "Error decoding GitHub user info:", err)
		return nil, err
	}

	// GitHub might not return an email, fetch it separately if needed
	if user.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if err := json.NewDecoder(emailResp.Body).Decode(&emails); err == nil {
				for _, e := range emails {
					if e.Primary {
						user.Email = e.Email
						break
					}
				}
			}
		}
	}

	o.logger.Debug(ctx, map[string]interface{}{
		"email": user.Email,
	}, "Get GitHub user info success")

	return &user, nil
}

func NewGithubOauth2Service(cfg *oauth2.Config, logger logger.Logger) Oauth2Service {
	return &GithubOauth2Service{
		oauth2Config: cfg,
		logger:       logger,
	}
}
