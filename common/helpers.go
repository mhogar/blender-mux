package common

import (
	"errors"
	"os"
)

func ChainError(message string, err error) error {
	return errors.New(message + "\n\t" + err.Error())
}

func GetEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "local" //default to local
	}

	return env
}
