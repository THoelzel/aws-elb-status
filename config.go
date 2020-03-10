package main

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type config struct {
}

func (c *config) GetPort() int {
	return viper.GetInt("port")
}

func NewConfig() *config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Error reading config file",
			zap.Error(err),
		)
	}

	viper.SetDefault("port", 8080)

	return &config{}
}
