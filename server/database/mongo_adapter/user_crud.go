package mongoadapter

import (
	"errors"

	"github.com/blendermux/common"
	"github.com/blendermux/server/models"
)

func (db *MongoAdapter) CreateUser(user *models.User) error {
	verr := user.Validate()
	if verr.Status != models.ModelValid {
		return common.ChainError("error validating user model", verr)
	}

	return errors.New("not implemented yet")
}

func (db *MongoAdapter) GetUserByEmail(email string) (*models.User, error) {
	return nil, errors.New("not implemented yet")
}
