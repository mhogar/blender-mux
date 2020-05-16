package config

import "github.com/spf13/viper"

type DatabaseConfig struct {
	URL     string
	Port    string
	Timeout int
	DBs     map[string]string
}

func initDatabaseConfig() {
	config := make(map[string]interface{})

	config["local"] = DatabaseConfig{
		URL:     "localhost",
		Port:    "27017",
		Timeout: 3000,
		DBs: map[string]string{
			"core":        "core",
			"integration": "integration",
		},
	}

	viper.Set("database", config)
}
