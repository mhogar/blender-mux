package migrationrunner

import (
	"errors"
	"log"

	"blendermux/common"
)

type Migration interface {
	GetTimestamp() string
	Up() error
	Down() error
}

type MigrationRepository interface {
	GetMigrations() []Migration
}

type MigrationCRUD interface {
	CreateMigration(timestamp string) error
	GetLatestTimestamp() (string, bool, error)
	DeleteMigrationByTimestamp(timestamp string) error
}

func MigrateUp(migrationRepo MigrationRepository, db MigrationCRUD) error {
	log.Println("Migrating Up")

	migrations := migrationRepo.GetMigrations()

	//get latest timestamp
	latestTimestamp, hasLatest, err := db.GetLatestTimestamp()
	if err != nil {
		return common.ChainError("error getting latest timestamp", err)
	}

	//print the timestamp if it exists
	if !hasLatest {
		log.Println("No timestamps found.")
	} else {
		log.Println("Latest timestamp:", latestTimestamp)
	}

	//run all migrations that are newer
	for _, migration := range migrations {
		timestamp := migration.GetTimestamp()

		if !hasLatest || timestamp > latestTimestamp {
			log.Println("Running", timestamp)

			err = migration.Up()
			if err != nil {
				return common.ChainError("error running migration", err)
			}

			//save the migration to db to mark it as run
			err = db.CreateMigration(timestamp)
			if err != nil {
				return common.ChainError("error saving migration", err)
			}
		} else {
			log.Println("Skipping", timestamp)
		}
	}

	log.Println("Finished running migrations.")
	return nil
}

func MigrateDown(migrationRepo MigrationRepository, db MigrationCRUD) error {
	log.Println("Migrating Down")

	migrations := migrationRepo.GetMigrations()

	//get latest timestamp
	latestTimestamp, hasLatest, err := db.GetLatestTimestamp()
	if err != nil {
		return common.ChainError("error getting latest timestamp", err)
	}

	//exit if no latest
	if !hasLatest {
		return errors.New("no migrations to migrate down")
	}

	var latestMigration Migration = nil

	//find migration that matches the latest timestamp
	for _, migration := range migrations {
		if migration.GetTimestamp() == latestTimestamp {
			latestMigration = migration
			break
		}
	}

	if latestMigration == nil {
		return errors.New("could not find migration with timestamp " + latestTimestamp)
	}

	log.Println("Running " + latestTimestamp)

	//run the down function
	err = latestMigration.Down()
	if err != nil {
		return common.ChainError("error running migration", err)
	}

	//remove migration from database
	err = db.DeleteMigrationByTimestamp(latestTimestamp)
	if err != nil {
		return common.ChainError("error deleting migration", err)
	}

	log.Println("Finished running migrations.")
	return nil
}
