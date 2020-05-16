package config

import "github.com/spf13/viper"

func InitConfig() {
	viper.SetDefault("env", "local")

	viper.SetEnvPrefix("cfg")
	viper.BindEnv("env")

	initDatabaseConfig()
}
