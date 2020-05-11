package mongomigrations

import (
	"context"

	"github.com/blendermux/common"

	mongoadapter "github.com/blendermux/server/database/mongo_adapter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type M20200507205301 struct {
	mongoadapter.MongoAdapter
}

func (m M20200507205301) GetTimestamp() string {
	return "20200507205301"
}

func (m M20200507205301) Up() error {
	//index users on email
	userIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"email", 1}},
		},
	}

	//set user indexes
	_, err := m.Users.Indexes().CreateMany(context.TODO(), userIndexes)
	if err != nil {
		return common.ChainError("error creating user indexes", err)
	}

	return nil
}

func (m M20200507205301) Down() error {
	//remove the created user indexes
	_, err := m.Users.Indexes().DropOne(context.TODO(), "email_1")
	if err != nil {
		return common.ChainError("error removing user indexes", err)
	}

	return nil
}
