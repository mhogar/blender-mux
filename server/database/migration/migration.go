package migration

type Migration interface {
	GetTimestamp() string
	Up() error
	Down() error
}

type MigrationRepository interface {
	GetMigrations() []Migration
}
