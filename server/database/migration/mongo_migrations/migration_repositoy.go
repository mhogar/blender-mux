package mongomigrations

import (
	"github.com/blendermux/server/database/migration"
	mongoadapter "github.com/blendermux/server/database/mongo_adapter"
)

type MongoMigrationRepository struct {
	mongoadapter.MongoAdapter
}

func (repo MongoMigrationRepository) GetMigrations() []migration.Migration {
	return []migration.Migration{
		M20200507205301{repo.MongoAdapter},
	}
}
