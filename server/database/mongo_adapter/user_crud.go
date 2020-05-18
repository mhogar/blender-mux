package mongoadapter

import (
	"errors"

	"blendermux/common"
	"blendermux/server/models"
)

func (db *MongoAdapter) CreateUser(user *models.User) error {
	verr := user.Validate()
	if verr.Status != models.ValidateErrorModelValid {
		return common.ChainError("error validating user model", verr)
	}

	return errors.New("not implemented yet")
}

func (db *MongoAdapter) GetUserByUsername(email string) (*models.User, error) {
	return nil, errors.New("not implemented yet")
}
