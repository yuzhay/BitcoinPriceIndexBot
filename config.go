package main

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Bitcoin  BitcoinConfig
	Telegram TelegramConfig
	DB       struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     uint
	}
}

type BitcoinConfig struct {
	Uri     string
	Timeout int
}

type TelegramConfig struct {
	Token string
}

func loadConfig(configPath string) (*Config, error) {
	var config Config

	fConfig, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("can't open YAML config %q: %s", configPath, err)
	}

	fData, err := ioutil.ReadAll(fConfig)
	if err != nil {
		return nil, fmt.Errorf("can't read YAML file %q: %s", configPath, err)
	}

	err = yaml.Unmarshal(fData, &config)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal YAML file %q: %s", configPath, err)
	}

	return &config, nil
}
