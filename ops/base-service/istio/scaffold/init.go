package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	appConfig = &Config{}
	if err = viper.Unmarshal(&appConfig); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
