package main

import (
	"log"

	"github.com/blendermux/server/database"
	"github.com/blendermux/server/database/migration"
	"github.com/blendermux/server/models"
)

func RunMigrations(migrationRepo migration.MigrationRepository, db database.MigrationCRUD) {
	//load the migrations
	migrations := migrationRepo.GetMigrations()

	//get latest timestamp
	latestTimestamp, err := db.GetLatestTimestamp()
	if err != nil {
		log.Fatal("Could not get latest timestamp:", err)
	}

	//print the timestamp if it exists
	if latestTimestamp == "" {
		log.Println("No timestamps found.")
	} else {
		log.Println("Latest timestamp:", latestTimestamp)
	}

	//run all migrations that are newer
	for _, migration := range migrations {
		timestamp := migration.GetTimestamp()

		if timestamp > latestTimestamp {
			log.Println("Running", timestamp)

			err = migration.Up()
			if err != nil {
				log.Fatal("Error running migration:", err)
			}

			//save the migration to db to mark it as run
			err = db.CreateNewMigration(models.CreateNewMigration(migration.GetTimestamp()))
			if err != nil {
				log.Fatal("Error saving migration:", err)
			}
		} else {
			log.Println("Skipping", timestamp)
		}
	}

	log.Println("Finished running migrations.")
}
