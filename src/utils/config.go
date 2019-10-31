package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// MongoConfig the config of mongo
type MongoConfig struct {
	URL               string `mapstructure:"url"`
	Port              string `mapstructure:"port"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	Database          string `mapstructure:"database"`
	ConnectionOptions string `mapstructure:"connectOptions"`
}

// Config the configuration of the app.
type Config struct {
	Mongo MongoConfig `mapstructure:"mongoConfig" yaml:"mongoConfig"`
}

var config Config

//GetConf return the config object
func GetConf() *Config {
	ex, _ := os.Getwd()
	appRootFolder := ex[:strings.Index(ex, "mooncake-duty")+14]
	viper.AddConfigPath(appRootFolder)
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
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
