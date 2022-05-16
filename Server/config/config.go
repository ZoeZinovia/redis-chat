package config

import (
	"encoding/json"
	"io/ioutil"
	logger "redisChat/Server/pkg/log"
)

var ServerConfig Config

type Config struct {
	MaxEntries int `json:"max_message_entries"`
}

// LoadConfig reads the config file and fills the ClientConfig global variable
func LoadConfig() (err error) {

	// Read config file
	configFile, err := ioutil.ReadFile("config/Config.json")
	if err != nil {
		logger.Logger.Error("reading config file", err, logger.Information{})
		return
	}

	// Unmarshal data into config struct
	if err = json.Unmarshal(configFile, &ServerConfig); err != nil {
		logger.Logger.Error("unmarshaling config data", err, logger.Information{})
		return
	}

	return
}
