package config

import (
	"errors"

	"github.com/blendermux/common"
)

type DatabaseConfig struct {
	URL     string            `json:"url"`
	Port    string            `json:"port"`
	Timeout int               `json:"timeout"`
	Dbs     map[string]string `json:"dbs"`
}

type databaseConfigMap map[string]DatabaseConfig

func (repo *ConfigFileRepository) GetDatabaseConfig() (DatabaseConfig, error) {
	//load the config if doesn't already exist
	if repo.databaseConfigMap == nil {
		repo.databaseConfigMap = make(databaseConfigMap)
		err := common.LoadConfigFromFile("config/database.json", &repo.databaseConfigMap)
		if err != nil {
			return DatabaseConfig{}, common.ChainError("error loading database config", err)
		}
	}

	env := common.GetEnv()
	config, ok := repo.databaseConfigMap[env]
	if !ok {
		return config, errors.New("no database config found for ENV " + env)
	}

	return config, nil
}
