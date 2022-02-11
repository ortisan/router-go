package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App           App           `mapstructure:"app"`
	Etcd          Etcd          `mapstructure:"etcd"`
	Redis         Redis         `mapstructure:"redis"`
	OpenTelemetry OpenTelemetry `mapstructure:"opentelemetry"`
	AWS           AWS           `mapstructure:"aws"`
}

type App struct {
	Name          string `mapstructure:"name"`
	ServerAddress string `mapstructure:"server_address"`
}

type Etcd struct {
	Endpoints []string `mapstructure:"server_endpoints"`
}

type Redis struct {
	ServerAddress string `mapstructure:"server_address"`
	Password      string `mapstructure:"password"`
}

type OpenTelemetry struct {
	AgentHost string `mapstructure:"agent_host"`
	AgentPort string `mapstructure:"agent_port"`
}

type SNS struct {
	HealthTopicArn string `mapstructure:"health_topic_arn"`
}

type SQS struct {
	HealthQueueUrl string `mapstructure:"health_queue_url"`
}

type AWS struct {
	Region      string `mapstructure:"region"`
	EndpointUrl string `mapstructure:"endpoint_url"`
	SNS         SNS    `mapstructure:"sns"`
	SQS         SQS    `mapstructure:"sqs"`
}

func LoadConfig() (config Config) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("../internal/config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	var err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return
}

var ConfigObj = LoadConfig()
