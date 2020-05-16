package dependencies

import (
	migrationrunner "blendermux/common/migration_runner"
	databasepkg "blendermux/server/database"

	mongoadapter "blendermux/server/database/mongo_adapter"
	"blendermux/server/database/mongo_adapter/migrations"
)

var database databasepkg.Database
var migrationRepository migrationrunner.MigrationRepository

func ResolveDatabase() databasepkg.Database {
	if database == nil {
		database = &mongoadapter.MongoAdapter{
			DbKey: "core",
		}
	}

	return database
}

func ResolveMigrationRepository() migrationrunner.MigrationRepository {
	if migrationRepository == nil {
		migrationRepository = &migrations.MongoMigrationRepository{
			MongoAdapter: ResolveDatabase().(*mongoadapter.MongoAdapter),
		}
	}

	return migrationRepository
}
