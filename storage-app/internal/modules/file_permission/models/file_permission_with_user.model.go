package models

type FilePermissionWithUser struct {
	FilePermission
	User *struct {
		Email     string `bson:"email"`
		FirstName string `bson:"first_name"`
		LastName  string `bson:"last_name"`
		Image     string `bson:"image"`
	}
}
