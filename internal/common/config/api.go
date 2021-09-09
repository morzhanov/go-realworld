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

type EventsServiceAPIItem struct {
	Event string `yaml:"event"`
}

type ServiceAPI struct {
	Rest   map[string]RestServiceAPIItem   `yaml:"rest"`
	Events map[string]EventsServiceAPIItem `yaml:"events"`
}

type ApiConfig struct {
	Analytics ServiceAPI `yaml:"analytics"`
	Auth      ServiceAPI `yaml:"auth"`
	Pictures  ServiceAPI `yaml:"pictures"`
	Users     ServiceAPI `yaml:"users"`
}

func NewApiConfig() (res *ApiConfig, err error) {
	data, err := os.ReadFile("./configs/api.yml")
	if err != nil {
		return nil, err
	}

	res = &ApiConfig{}
	if err = yaml.Unmarshal([]byte(data), res); err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ApiConfig) GetApiItem(key string) (*ServiceAPI, error) {
	switch key {
	case "analytics":
		return &a.Analytics, nil
	case "auth":
		return &a.Auth, nil
	case "pictures":
		return &a.Pictures, nil
	case "users":
		return &a.Users, nil
	default:
		return nil, fmt.Errorf("wrong api key %v", key)
	}
}
