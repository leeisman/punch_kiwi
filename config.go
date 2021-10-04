package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"path/filepath"
)

var _config *Config

type User struct {
	Username string
	ID       string
	Password string
}

// Config ...
type Config struct {
	Users []*User
}

// NewInjection ...
func (c *Config) NewInjection() *Config {
	return c
}

func NewConfig() (*Config, error) {
	dir, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}
	fmt.Println(dir)
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)
	viper.SetConfigType("yaml")

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Msgf("Error reading config file, %s", err)
		return &config, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
