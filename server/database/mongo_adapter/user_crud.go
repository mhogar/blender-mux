package mongoadapter

import (
	"errors"

	"blendermux/common"
	"blendermux/server/models"

	"github.com/google/uuid"
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

func (db *MongoAdapter) GetUserBySessionId(sID uuid.UUID) (*models.User, error) {
	return nil, errors.New("not implemented yet")
}

func (db *MongoAdapter) UpdateUser(user *models.User) error {
	return errors.New("not implemented yet")
}

func (db *MongoAdapter) DeleteUser(id uuid.UUID) (bool, error) {
	return false, errors.New("not implemented yet")
}
