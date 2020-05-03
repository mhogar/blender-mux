package inmemorydb

import (
	"github.com/blendermux/server/models"
	"github.com/google/uuid"
)

func (db InMemoryDB) CreateSession(session *models.Session) {
	db.Sessions[session.ID] = session
}

func (db InMemoryDB) GetSessionByID(id uuid.UUID) *models.Session {
	return db.Sessions[id]
}

func (db InMemoryDB) DeleteSession(session *models.Session) {
	delete(db.Sessions, session.ID)
}
