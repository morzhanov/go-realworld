package config

import "github.com/spf13/viper"

type Config struct {
	KafkaUri string `mapstructure:"KAFKA_URI"`

	AnalyticsKafkaTopic string `mapstructure:"ANALYTICS_KAFKA_TOPIC"`
	AuthKafkaTopic      string `mapstructure:"AUTH_KAFKA_TOPIC"`
	PicturesKafkaTopic  string `mapstructure:"PICTURES_KAFKA_TOPIC"`
	UsersKafkaTopic     string `mapstructure:"USERS_KAFKA_TOPIC"`
	ResultsKafkaTopic   string `mapstructure:"RESULTS_KAFKA_TOPIC"`

	PsqlConnectionString string `mapstructure:"PSQL_CONNECTION_STRING"`

	AccessTokenSecret string `mapstructure:"ACCESS_TOKEN_SECRET"`

	AnalyticsGrpcAddr string `mapstructure:"ANALYTICS_GRPC_ADDR"`
	AuthGrpcAddr      string `mapstructure:"AUTH_GRPC_ADDR"`
	PicturesGrpcAddr  string `mapstructure:"PICTURES_GRPC_ADDR"`
	UsersGrpcAddr     string `mapstructure:"USERS_GRPC_ADDR"`

	AnalyticsGrpcPort string `mapstructure:"ANALYTICS_GRPC_PORT"`
	AuthGrpcPort      string `mapstructure:"AUTH_GRPC_PORT"`
	PicturesGrpcPort  string `mapstructure:"PICTURES_GRPC_PORT"`
	UsersGrpcPort     string `mapstructure:"USERS_GRPC_PORT"`

	AnalyticsRestPort string `mapstructure:"ANALYTICS_REST_PORT"`
	AuthRestPort      string `mapstructure:"AUTH_REST_PORT"`
	PicturesRestPort  string `mapstructure:"PICTURES_REST_PORT"`
	UsersRestPort     string `mapstructure:"USERS_REST_PORT"`
}

func LoadConfig(path string) (config Config, err error) {
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
