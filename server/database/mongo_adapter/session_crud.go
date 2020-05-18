package mongoadapter

import (
	"errors"

	"blendermux/common"
	"blendermux/server/models"

	"github.com/google/uuid"
)

func (db *MongoAdapter) CreateSession(session *models.Session) error {
	verr := session.Validate()
	if verr.Status != models.ValidateErrorModelValid {
		return common.ChainError("error validating session model", verr)
	}

	return errors.New("not implemented yet")
}

func (db *MongoAdapter) GetSessionByID(ID uuid.UUID) (*models.Session, error) {
	return nil, errors.New("not implemented yet")
}

func (db *MongoAdapter) DeleteSession(session *models.Session) error {
	return errors.New("not implemented yet")
}
