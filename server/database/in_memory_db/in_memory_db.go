package inmemorydb

import (
	"github.com/blendermux/server/models"
	"github.com/google/uuid"
)

type InMemoryDB struct {
	Users    map[uuid.UUID]*models.User
	Sessions map[uuid.UUID]*models.Session
}

func InitInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		make(map[uuid.UUID]*models.User),
		make(map[uuid.UUID]*models.Session),
	}
}
