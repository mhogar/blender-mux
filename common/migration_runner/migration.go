package migrationrunner

// Migration is an interface for running a migration.
type Migration interface {
	// GetTimestamp should return the timestamp for this migration.
	GetTimestamp() string

	// Up runs the migration and returns any errors.
	Up() error

	// Down runs the inverse migration and returns any errors.
	Down() error
}

// MigrationRepository is an interface for fetching migrations to run.
type MigrationRepository interface {
	// GetMigrations returns the migrations to run
	GetMigrations() []Migration
}
