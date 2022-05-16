package config

import (
	"encoding/json"
	"io/ioutil"
	logger "redisChat/Client/pkg/log"
)

var ClientConfig Config
var User string
var ClientPort string

type Config struct {
	ServerAddressReceiver string `json:"server_address_receiver"`
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
	if err = json.Unmarshal(configFile, &ClientConfig); err != nil {
		logger.Logger.Error("unmarshaling config data", err, logger.Information{})
		return
	}

	return
}
