package config

import (
	"StoreServer/utils/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var config *Config

type Config struct {
	ServerPort string `envconfig:"SERVER_PORT"`
	MONGOURL   string `envconfig:"MONGO_URL"`
}

func init() {
	GetConfig()
}

func GetConfig() *Config {
	config = &Config{}

	_ = godotenv.Load()

	err := envconfig.Process("", config)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	return config
}
