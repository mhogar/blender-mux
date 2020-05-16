package migrations

import (
	"blendermux/common"

	mongoadapter "blendermux/server/database/mongo_adapter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type M20200507205301 struct {
	*mongoadapter.MongoAdapter
}

func (m M20200507205301) GetTimestamp() string {
	return "20200507205301"
}

func (m M20200507205301) Up() error {
	//index users on email
	userIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
		},
	}

	//set user indexes
	ctx, cancel := m.CreateStandardTimeoutContext()
	_, err := m.Users.Indexes().CreateMany(ctx, userIndexes)
	cancel()

	if err != nil {
		return common.ChainError("error creating user indexes", err)
	}

	return nil
}

func (m M20200507205301) Down() error {
	//remove the created user indexes
	ctx, cancel := m.CreateStandardTimeoutContext()
	_, err := m.Users.Indexes().DropOne(ctx, "email_1")
	cancel()

	if err != nil {
		return common.ChainError("error removing user indexes", err)
	}

	return nil
}
