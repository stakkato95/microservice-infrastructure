package config

import "github.com/stakkato95/service-engineering-go-lib/config"

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	DbSource   string `mapstructure:"DB_SOURCE"`
}

var AppConfig Config

func init() {
	config.Init(&AppConfig, Config{})
}
