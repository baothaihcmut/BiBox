package users

type UserSignUpEvent struct {
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	ConfirmationLink string `json:"confirmation_link"`
}
