package common

import (
	"encoding/json"
	"os"
)

func LoadConfigFromFile(filename string, config interface{}) error {
	//open the file
	file, err := os.Open(filename)
	if err != nil {
		return ChainError("error opening config file", err)
	}
	defer file.Close()

	//decode the json from the file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return ChainError("error parsing config json", err)
	}

	return nil
}
