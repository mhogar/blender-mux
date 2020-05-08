package main

import (
	"log"

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

	RunMigrations(resolver.MigrationRepository, db)
}
