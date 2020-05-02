package database

import "github.com/blendermux/server/models"

type InMemoryDB struct {
	Users map[string]*models.User
}

func InitInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		make(map[string]*models.User),
	}
}

func (db InMemoryDB) CreateUser(user *models.User) {
	db.Users[user.Email] = user
}

func (db InMemoryDB) GetUserByEmail(email string) *models.User {
	return db.Users[email]
}
