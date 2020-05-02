package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User foo
type User struct {
	ID           primitive.ObjectID
	Email        string
	PasswordHash []byte
}
