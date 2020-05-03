package inmemorydb

import (
	"github.com/blendermux/server/models"
)

func (db InMemoryDB) CreateUser(user *models.User) {
	db.Users[user.ID] = user
}

func (db InMemoryDB) GetUserByEmail(email string) *models.User {
	for _, value := range db.Users {
		if value.Email == email {
			return value
		}
	}

	return nil
}
