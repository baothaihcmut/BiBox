package presenter

type SignUpInput struct {
	Email          string `json:"email" bindinng:"required"`
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeat_password" binding:"required"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
}

type SignUpOutput struct{}
