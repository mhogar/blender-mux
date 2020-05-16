package main

import (
	"flag"
	"log"

	"blendermux/server/config"

	"blendermux/common"

	migrationrunner "blendermux/common/migration_runner"
	"blendermux/server/dependencies"
)

func main() {
	config.InitConfig()

	//parse flags
	down := flag.Bool("down", false, "Run the most recent migration down")
	flag.Parse()

	db := dependencies.ResolveDatabase()

	//open the db connection
	err := db.OpenConnection()
	if err != nil {
		log.Fatal(common.ChainError("Could not create database connection", err))
	}

	defer db.CloseConnection()

	//check db is connected
	err = db.Ping()
	if err != nil {
		log.Fatal(common.ChainError("Could not reach database", err))
	}

	//run the migrations
	if *down {
		err = migrationrunner.MigrateDown(dependencies.ResolveMigrationRepository(), db)
	} else {
		err = migrationrunner.MigrateUp(dependencies.ResolveMigrationRepository(), db)
	}

	if err != nil {
		log.Fatal(common.ChainError("Error running migrations", err))
	}
}
