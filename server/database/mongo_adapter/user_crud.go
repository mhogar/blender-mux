package mongoadapter

import (
	"errors"

	"github.com/blendermux/server/models"
)

func (db MongoAdapter) CreateUser(user *models.User) error {
	return errors.New("Not implemented yet.")
}

func (db MongoAdapter) GetUserByEmail(email string) (*models.User, error) {
	return nil, errors.New("Not implemented yet.")
}
