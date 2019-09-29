package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

// MongoConfig the config of mongo
type MongoConfig struct {
	URL               string `json:"url"`
	Port              string `json:"port"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Database          string `json:"database"`
	ConnectionOptions string `json:"connectOptions"`
}

// Config the configuration of the app.
type Config struct {
	Mongo MongoConfig `json:"mongoConfig"`
}

var config Config

//GetConf return the config object
func GetConf() *Config {
	viper.SetConfigName("config")               // name of config file (without extension)
	viper.AddConfigPath("/etc/mooncake-duty/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.mooncake-duty") // call multiple times to add many search paths
	viper.AddConfigPath(".")                    // optionally look for config in the working directory
	err := viper.ReadInConfig()                 // Find and read the config file
	if err != nil {                             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	return &config
}
