package migrationrunner

import (
	"log"

	"github.com/blendermux/common"
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
}

func RunMigrations(migrationRepo MigrationRepository, db MigrationCRUD) error {
	//load the migrations
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
