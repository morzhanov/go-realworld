package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type RestServiceAPIItem struct {
	Method string `yaml:"method"`
	Url    string `yaml:"url"`
}

type GrpcServiceAPIItem struct {
	Method string `yaml:"method"`
}

type EventsServiceAPIItem struct {
	Event string `yaml:"event"`
}

type ServiceAPI struct {
	Rest   map[string]RestServiceAPIItem   `yaml:"rest"`
	Grpc   map[string]GrpcServiceAPIItem   `yaml:"grpc"`
	Events map[string]EventsServiceAPIItem `yaml:"events"`
}

type apiConfig struct {
	analytics ServiceAPI `yaml:"analytics"`
	auth      ServiceAPI `yaml:"auth"`
	pictures  ServiceAPI `yaml:"pictures"`
	users     ServiceAPI `yaml:"users"`
}

type ApiConfig interface {
	GetApiItem(key string) (*ServiceAPI, error)
	Analytics() ServiceAPI
	Auth() ServiceAPI
	Pictures() ServiceAPI
	Users() ServiceAPI
}

func (a *apiConfig) GetApiItem(key string) (*ServiceAPI, error) {
	switch key {
	case "analytics":
		return &a.analytics, nil
	case "auth":
		return &a.auth, nil
	case "pictures":
		return &a.pictures, nil
	case "users":
		return &a.users, nil
	default:
		return nil, fmt.Errorf("wrong api key %v", key)
	}
}

func (a *apiConfig) Analytics() ServiceAPI { return a.analytics }
func (a *apiConfig) Auth() ServiceAPI      { return a.auth }
func (a *apiConfig) Pictures() ServiceAPI  { return a.pictures }
func (a *apiConfig) Users() ServiceAPI     { return a.users }

func NewApiConfig() (res ApiConfig, err error) {
	data, err := os.ReadFile("./configs/api.yml")
	if err != nil {
		return nil, err
	}
	res = &apiConfig{}
	if err = yaml.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}
