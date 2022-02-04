package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App   App   `mapstructure:"app"`
	Etcd  Etcd  `mapstructure:"etcd"`
	Redis Redis `mapstructure:"redis"`
}

type App struct {
	ServerAddress string `mapstructure:"server_address"`
}

type Etcd struct {
	Endpoints []string `mapstructure:"server_endpoints"`
}

type Redis struct {
	ServerAddress string `mapstructure:"server_address"`
	Password      string `mapstructure:"password"`
}

func LoadConfig() (config Config) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("../internal/config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err2 := viper.Unmarshal(&config)
	if err != nil {
		panic(err2)
	}

	return
}

var ConfigObj = LoadConfig()
