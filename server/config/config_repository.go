package config

type ConfigRepository interface {
	GetDatabaseConfig() (DatabaseConfig, error)
}

type ConfigFileRepository struct {
	databaseConfigMap databaseConfigMap
}
