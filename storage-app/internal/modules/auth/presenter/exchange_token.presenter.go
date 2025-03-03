package presenter

type ExchangeTokenInput struct {
	AuthCode string `json:"auth_code"`
	Provider int    `json:"provider" binding:"min=1,max=2"`
}

type ExchangeTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
