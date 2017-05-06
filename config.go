package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// NewConfig initializes configuration to use N26 API
func NewConfig() *N26API {
	config := viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("n26")
	config.AddConfigPath("$HOME/.config")
	err := config.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read config, %s", err)
		return nil
	}
	return &N26API{config.GetString("username"), config.GetString("password")}
}
