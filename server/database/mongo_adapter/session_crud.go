package mongoadapter

import (
	"errors"

	"github.com/blendermux/server/models"
	"github.com/google/uuid"
)

func (db MongoAdapter) CreateSession(session *models.Session) error {
	verr := session.Validate()
	if verr.Status != models.ModelValid {
		return errors.New("error validating session model: " + verr.Error())
	}

	return errors.New("not implemented yet")
}

func (db MongoAdapter) GetSessionByToken(token uuid.UUID) (*models.Session, error) {
	return nil, errors.New("not implemented yet")
}

func (db MongoAdapter) DeleteSession(session *models.Session) error {
	return errors.New("not implemented yet")
}
