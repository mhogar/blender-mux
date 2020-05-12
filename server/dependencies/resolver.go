package dependencies

import (
	migrationrunner "github.com/blendermux/common/migration_runner"

	"github.com/blendermux/server/database"
	mongoadapter "github.com/blendermux/server/database/mongo_adapter"
	"github.com/blendermux/server/database/mongo_adapter/migrations"
)

type DependencyResolver struct {
	database.Database
	migrationrunner.MigrationRepository
}

func CreateDependencyResolver() DependencyResolver {
	database := &mongoadapter.MongoAdapter{}
	migrationRepo := &migrations.MongoMigrationRepository{MongoAdapter: database}

	return DependencyResolver{
		Database:            database,
		MigrationRepository: migrationRepo,
	}
}
