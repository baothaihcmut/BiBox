package presenter

type ExchangeTokenInput struct {
	AuthCode string `json:"auth_code"`
}

type ExchangeTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
