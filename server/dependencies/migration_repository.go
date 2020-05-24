package dependencies

import (
	migrationrunner "blendermux/common/migration_runner"
	mongoadapter "blendermux/server/database/mongo_adapter"
	"blendermux/server/database/mongo_adapter/migrations"
	"sync"
)

var createMigrationRepositoryOnce sync.Once
var migrationRepository migrationrunner.MigrationRepository

// ResolveMigrationRepository resolves the MigrationRepository dependency.
// Only the first call to this function will create a new MigrationRepository, after which it will be retrieved from the cache.
func ResolveMigrationRepository() migrationrunner.MigrationRepository {
	createMigrationRepositoryOnce.Do(func() {
		migrationRepository = migrations.MongoMigrationRepository{
			MongoAdapter: ResolveDatabase().(*mongoadapter.MongoAdapter),
		}
	})
	return migrationRepository
}
