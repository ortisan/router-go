package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

type Config struct {
	App           App           `mapstructure:"app"`
	Etcd          Etcd          `mapstructure:"etcd"`
	Redis         Redis         `mapstructure:"redis"`
	OpenTelemetry OpenTelemetry `mapstructure:"opentelemetry"`
	AWS           AWS           `mapstructure:"aws"`
	Servers       []Server      `mapstructure:"servers"`
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

type Credentials struct {
	AccessKeyID  string `mapstructure:"access_key_id"`
	ClientSecret string `mapstructure:"secret_access_key"`
}

type AWS struct {
	Region      string      `mapstructure:"region"`
	Credentials Credentials `mapstructure:"credentials"`
	EndpointUrl string      `mapstructure:"endpoint_url"`
	SNS         SNS         `mapstructure:"sns"`
	SQS         SQS         `mapstructure:"sqs"`
}

type HealthCheck struct {
	Type     string `mapstructure:"type"`
	Endpoint string `mapstructure:"endpoint"`
}

type Server struct {
	ServicePrefix string      `mapstructure:"service_prefix"`
	ServerName    string      `mapstructure:"server_name"`
	EndpointUrl   string      `mapstructure:"endpoint_url"`
	ZoneAws       string      `mapstructure:"zone_aws"`
	Alive         bool        `mapstructure:"alive"`
	HealthCheck   HealthCheck `mapstructure:"healthcheck"`
}

func Setup() (config Config) {
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

func SetupAws() *aws.Config {
	awsConfig := aws.NewConfig()
	// Config credentials
	if ConfigObj.AWS.Credentials.AccessKeyID != "" {
		awsConfig.WithCredentials(credentials.NewStaticCredentials(ConfigObj.AWS.Credentials.AccessKeyID, ConfigObj.AWS.Credentials.ClientSecret, ""))
	}
	// Config region
	awsConfig.WithRegion(ConfigObj.AWS.Region)
	awsConfig.WithS3ForcePathStyle(true)

	// Config endpoint url (local and docker env)
	if len(ConfigObj.AWS.EndpointUrl) > 0 {
		awsConfig.WithEndpoint(ConfigObj.AWS.EndpointUrl)
	}

	return awsConfig
}

func NewAWSSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *AwsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))
}

var ConfigObj = Setup()
var AwsConfig = SetupAws()
