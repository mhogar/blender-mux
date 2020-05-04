package inmemorydb

import (
	"github.com/blendermux/server/models"
)

func (db InMemoryDB) CreateUser(user *models.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	db.Users[user.ID] = user
	return nil
}

func (db InMemoryDB) GetUserByEmail(email string) *models.User {
	for _, value := range db.Users {
		if value.Email == email {
			return value
		}
	}

	return nil
}
