package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Config struct {
	OutputFolder string `json:"outputFolder"`
	FilePrefix   string `json:"filePrefix"`
}

func loadConfig(jsonConfig io.Reader) (Config, error) {

	readConfig, readAllError := ioutil.ReadAll(jsonConfig)
	if readAllError != nil {
		return defaultConfig, readAllError
	}

	unmarshalError := json.Unmarshal(readConfig, &defaultConfig)
	if unmarshalError != nil {
		return defaultConfig, unmarshalError
	}

	return defaultConfig, nil
}
