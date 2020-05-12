package dependencies

import (
	mongomigrations "github.com/blendermux/server/database/migrations/mongo_migrations"

	migrationrunner "github.com/blendermux/common/migration_runner"

	"github.com/blendermux/server/database"
	mongoadapter "github.com/blendermux/server/database/mongo_adapter"
)

type DependencyResolver struct {
	database.Database
	migrationrunner.MigrationRepository
}

func GetDependencyResolver() DependencyResolver {
	database := &mongoadapter.MongoAdapter{}
	migrationRepo := &mongomigrations.MongoMigrationRepository{database}

	return DependencyResolver{
		database,
		migrationRepo,
	}
}
