package common

import "errors"

func ChainError(message string, err error) error {
	return errors.New(message + ": " + err.Error())
}
