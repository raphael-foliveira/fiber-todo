package common

import (
	"encoding/json"
	"os"
)

type DatabaseConfig struct {
	Url string `json:"url"`
}

type AppConfig struct {
	Port int `json:"port"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	App      AppConfig      `json:"app"`
}

func ReadTestCfg() Config {
	config := Config{}
	cfgB, err := os.ReadFile("../../config_test.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(cfgB, &config)
	if err != nil {
		panic(err)
	}
	return config
}
