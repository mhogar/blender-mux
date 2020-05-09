package mongoadapter

import (
	"errors"

	"github.com/blendermux/server/models"
	"github.com/google/uuid"
)

func (db MongoAdapter) CreateSession(session *models.Session) error {
	return errors.New("Not implemented yet.")
}

func (db MongoAdapter) GetSessionByID(id uuid.UUID) (*models.Session, error) {
	return nil, errors.New("Not implemented yet.")
}

func (db MongoAdapter) DeleteSession(session *models.Session) error {
	return errors.New("Not implemented yet.")
}
