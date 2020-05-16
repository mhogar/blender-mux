package migrations

import (
	migrationrunner "blendermux/common/migration_runner"
	mongoadapter "blendermux/server/database/mongo_adapter"
)

// MongoMigrationRepository is a struct with the pointer to the mongo adapter that its migrations will be run against.
type MongoMigrationRepository struct {
	*mongoadapter.MongoAdapter
}

// GetMigrations returns a slice of Migrations that need to be run on the MongoDB database.
func (repo MongoMigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{}
}
