package dependencies

import (
	migrationrunner "github.com/blendermux/common/migration_runner"
	"github.com/blendermux/server/config"

	"github.com/blendermux/server/database"
	mongoadapter "github.com/blendermux/server/database/mongo_adapter"
	"github.com/blendermux/server/database/mongo_adapter/migrations"
)

type DependencyResolver struct {
	config.ConfigRepository
	database.Database
	migrationrunner.MigrationRepository
}

func CreateDependencyResolver() DependencyResolver {
	configRepo := &config.ConfigFileRepository{}
	database := &mongoadapter.MongoAdapter{ConfigRepository: configRepo, DbKey: "core"}
	migrationRepo := &migrations.MongoMigrationRepository{MongoAdapter: database}

	return DependencyResolver{
		ConfigRepository:    configRepo,
		Database:            database,
		MigrationRepository: migrationRepo,
	}
}
