package config

import "github.com/spf13/viper"

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`

	KafkaUri                    string `mapstructure:"KAFKA_URI"`
	KafkaTopic                  string `mapstructure:"KAFKA_TOPIC"`
	AnalyticsKafkaTopic         string `mapstructure:"ANALYTICS_KAFKA_TOPIC"`
	AuthKafkaTopic              string `mapstructure:"AUTH_KAFKA_TOPIC"`
	PicturesKafkaTopic          string `mapstructure:"PICTURES_KAFKA_TOPIC"`
	UsersKafkaTopic             string `mapstructure:"USERS_KAFKA_TOPIC"`
	ResultsKafkaTopic           string `mapstructure:"RESULTS_KAFKA_TOPIC"`
	KafkaAnalyticsDbTopic       string `mapstructure:"KAFKA_ANALYTICS_DB_TOPIC"`
	KafkaConsumerGroupId        string `mapstructure:"KAFKA_CONSUMER_GROUP_ID"`
	KafkaResultsConsumerGroupId string `mapstructure:"KAFKA_RESULTS_CONSUMER_GROUP_ID"`

	PsqlConnectionString string `mapstructure:"PSQL_CONNECTION_STRING"`

	AccessTokenSecret string `mapstructure:"ACCESS_TOKEN_SECRET"`

	GrpcAddr          string `mapstructure:"GRPC_ADDR"`
	GrpcPort          string `mapstructure:"GRPC_PORT"`
	PicturesGrpcPort  string `mapstructure:"PICTURES_GRPC_PORT"`
	UsersGrpcPort     string `mapstructure:"USERS_GRPC_PORT"`
	AnalyticsGrpcPort string `mapstructure:"ANALYTICS_GRPC_PORT"`
	AuthGrpcPort      string `mapstructure:"AUTH_GRPC_PORT"`

	RestAddr          string `mapstructure:"REST_ADDR"`
	RestPort          string `mapstructure:"REST_PORT"`
	AnalyticsRestPort string `mapstructure:"ANALYTICS_REST_PORT"`
	AuthRestPort      string `mapstructure:"AUTH_REST_PORT"`
	PicturesRestPort  string `mapstructure:"PICTURES_REST_PORT"`
	UsersRestPort     string `mapstructure:"USERS_REST_PORT"`
	ApiGWRestPort     string `mapstructure:"APIGW_REST_PORT"`
}

func NewConfig(path string, name string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
