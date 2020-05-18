package dependencies

import (
	migrationrunner "blendermux/common/migration_runner"
	"blendermux/server/controllers"
	databasepkg "blendermux/server/database"
	mongoadapter "blendermux/server/database/mongo_adapter"
	"blendermux/server/database/mongo_adapter/migrations"
)

var database databasepkg.Database
var migrationRepository migrationrunner.MigrationRepository
var requestHandler controllers.RequestHandler

// ResolveDatabase resolves the Database dependency.
// Only the first call to this function will create a new Database, after which it will be retrieved from the cache.
func ResolveDatabase() databasepkg.Database {
	if database == nil {
		database = &mongoadapter.MongoAdapter{
			DbKey: "core",
		}
	}

	return database
}

// ResolveMigrationRepository resolves the MigrationRepository dependency.
// Only the first call to this function will create a new MigrationRepository, after which it will be retrieved from the cache.
func ResolveMigrationRepository() migrationrunner.MigrationRepository {
	if migrationRepository == nil {
		migrationRepository = &migrations.MongoMigrationRepository{
			MongoAdapter: ResolveDatabase().(*mongoadapter.MongoAdapter),
		}
	}

	return migrationRepository
}

// ResolveRouteHandler resolves the RouteHandler dependency.
// Only the first call to this function will create a new RouteHandler, after which it will be retrieved from the cache.
func ResolveRouteHandler() controllers.RequestHandler {
	if requestHandler == nil {
		requestHandler = &controllers.RequestHandle{
			UserController: controllers.UserController{
				UserCRUD: ResolveDatabase(),
			},
		}
	}

	return requestHandler
}
