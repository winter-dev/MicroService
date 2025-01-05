package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var G_AppConfig GlobalConfig

type GlobalConfig struct {
	AppConfig AppConfig
}

type AppConfig struct {
	AppName    string   `mapstructure:"appName"`
	Version    string   `mapstructure:"version"`
	Env        string   `mapstructure:"env"`
	Address    string   `mapstructure:"address"`
	GrpcPort   int      `mapstructure:"grpcPort"`
	HttpPort   int      `mapstructure:"httpPort"`
	EcdAddress []string `mapstructure:"etcdAddress"`
}

func Init() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config load error : %s", err.Error()))
	}
	if err := viper.Unmarshal(&G_AppConfig); err != nil {
		panic(fmt.Errorf("config unmarshal error : %s", err.Error()))
	}
}
