package dependencies

import (
	"log"

	"github.com/blendermux/common"
	mongomigrations "github.com/blendermux/server/database/migrations/mongo_migrations"

	migrationrunner "github.com/blendermux/common/migration_runner"

	"github.com/blendermux/server/database"
	mongoadapter "github.com/blendermux/server/database/mongo_adapter"
)

type DependencyResolver struct {
	database.Database
	migrationrunner.MigrationRepository
}

func CreateDependencyResolver() *DependencyResolver {
	//init database dependency
	database := mongoadapter.MongoAdapter{}
	err := database.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	//init migration repository dependency
	migrationRepo := mongomigrations.MongoMigrationRepository{database}

	return &DependencyResolver{
		database,
		migrationRepo,
	}
}

func (resolver DependencyResolver) DestroyDependencies() error {
	err := resolver.Database.Destroy()
	if err != nil {
		return common.ChainError("error destorying database dependency", err)
	}

	return nil
}
