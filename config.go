package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// NewConfig initializes the config file
func NewConfig(username, password string) *N26Credentials {
	return &N26Credentials{username, password}
}

// Config returns configuration from file to use N26 API
func Config() *N26Credentials {
	config := viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("n26")
	config.AddConfigPath("$HOME/.config")
	err := config.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read config, %s", err)
		return nil
	}
	return &N26Credentials{config.GetString("username"), config.GetString("password")}
}
