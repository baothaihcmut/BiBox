package users

type UserSignUpEvent struct {
	Email            string
	FirstName        string
	LastName         string
	ConfirmationLink string
}
