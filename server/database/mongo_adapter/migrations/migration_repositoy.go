package migrations

import (
	migrationrunner "blendermux/common/migration_runner"
	mongoadapter "blendermux/server/database/mongo_adapter"
)

type MongoMigrationRepository struct {
	*mongoadapter.MongoAdapter
}

func (repo MongoMigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{
		M20200507205301{MongoAdapter: repo.MongoAdapter},
	}
}
