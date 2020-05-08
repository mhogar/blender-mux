package main

import (
	"log"

	migrationrunner "github.com/blendermux/common/migration_runner"
	"github.com/blendermux/server/dependencies"
)

func main() {
	resolver := dependencies.CreateDependencyResolver()
	defer resolver.DestroyDependencies()

	db := resolver.Database

	//check db connection
	err := db.Ping()
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	migrationrunner.RunMigrations(resolver.MigrationRepository, db)
}
