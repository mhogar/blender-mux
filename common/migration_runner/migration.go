package migrationrunner

// Migration is an interface for running a migration.
//
// GetTimestamp should return the timestamp for this migration.
//
// Up runs the migration and returns any errors.
//
// Down runs the inverse migration and returns any errors.
type Migration interface {
	GetTimestamp() string
	Up() error
	Down() error
}

// MigrationRepository is an interface for fetching migrations to run.
type MigrationRepository interface {
	GetMigrations() []Migration
}
