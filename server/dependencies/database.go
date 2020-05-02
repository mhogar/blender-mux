package dependencies

import "github.com/blendermux/server/models"

type Database interface {
	UserCRUD
}

type UserCRUD interface {
	CreateUser(user *models.User)
	GetUserByEmail(email string) *models.User
}
