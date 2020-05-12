package migrations

import (
	migrationrunner "github.com/blendermux/common/migration_runner"
	mongoadapter "github.com/blendermux/server/database/mongo_adapter"
)

type MongoMigrationRepository struct {
	*mongoadapter.MongoAdapter
}

func (repo MongoMigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{
		M20200507205301{MongoAdapter: repo.MongoAdapter},
	}
}
