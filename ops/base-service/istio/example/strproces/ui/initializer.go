package main

import (
	"github.com/spf13/viper"
)

var (
	viperInst *viper.Viper
	err       error
)

func init() {
	initConfig()

}

func initConfig() {
	if appConfig == nil {
		appConfig = &Config{}
	}

	viperInst = viper.New()
	viperInst.SetConfigName("config")
	viperInst.AddConfigPath(".")
	viperInst.SetConfigType("json")

	err := viperInst.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viperInst.Unmarshal(appConfig)
}
