package models

type FilePermissionWithUser struct {
	*FilePermission `bson:"inline"`
	User            *struct {
		Email     string `bson:"email"`
		FirstName string `bson:"first_name"`
		LastName  string `bson:"last_name"`
		Image     string `bson:"image"`
	}
}
