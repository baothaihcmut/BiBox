package presenter

type ExchangeTokenInput struct {
	AuthCode string
}

type ExchangeTokenOutput struct {
	AccessToken  string
	RefreshToken string
}
