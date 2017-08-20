package main

import "log"

var config *Config

func init() {
	var err error

	config, err = loadConfig("config.yml")
	if err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	err = dbConnect(config.DB.Name, config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port)
	if err != nil {
		log.Fatalf("can't connect to database: %s", err)
	}
}
