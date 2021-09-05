package config

import "github.com/spf13/viper"

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`

	KafkaUri   string `mapstructure:"KAFKA_URI"`
	KafkaTopic string `mapstructure:"KAFKA_TOPIC"`

	PsqlConnectionString string `mapstructure:"PSQL_CONNECTION_STRING"`

	AccessTokenSecret string `mapstructure:"ACCESS_TOKEN_SECRET"`

	GrpcAddr string `mapstructure:"GRPC_ADDR"`
	GrpcPort string `mapstructure:"GRPC_PORT"`

	RestPort string `mapstructure:"REST_PORT"`
}

func NewConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
