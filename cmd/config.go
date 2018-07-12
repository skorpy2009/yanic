package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"chaos.expert/FreifunkBremen/yanic/database"
	"chaos.expert/FreifunkBremen/yanic/respond"
	"chaos.expert/FreifunkBremen/yanic/runtime"
	"chaos.expert/FreifunkBremen/yanic/webserver"
	"github.com/naoina/toml"
)

// Config represents the whole configuration
type Config struct {
	Respondd  respond.Config
	Webserver webserver.Config
	Nodes     runtime.NodesConfig
	Database  database.Config
}

var (
	configPath string
	collector  *respond.Collector
	nodes      *runtime.Nodes
)

func loadConfig() *Config {
	config, err := ReadConfigFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to load config file:", err)
		os.Exit(2)
	}
	return config
}

// ReadConfigFile reads a config model from path of a yml file
func ReadConfigFile(path string) (config *Config, err error) {
	config = &Config{}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return
}
