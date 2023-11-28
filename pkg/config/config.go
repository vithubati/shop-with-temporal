package config

import (
	"github.com/jinzhu/configor"
)

var (
	Configor *configor.Configor
)

type Configuration struct {
	Services *Services
	Env      string
	Temporal Temporal
}
type Services struct {
	Shopper *Service
}
type Service struct {
	APPName       string `default:"app name"`
	Host          string
	Port          int
	Version       string
	DB            DB
	LogLevel      uint32
	LogTimeFormat string
}
type Temporal struct {
	HostPort string
}

type DB struct {
	Dialect  int32  `default:"1"`
	Database string `default:"shopper.db"`
}

func NewConfig(file string) (*Configuration, error) {
	config := &Configuration{}

	Configor = configor.New(&configor.Config{Debug: true, ErrorOnUnmatchedKeys: true})
	if err := Configor.Load(config, file); err != nil {
		return nil, err
	}
	config.Env = Configor.GetEnvironment()
	return config, nil
}

func (c Configuration) IsDev() bool {
	return c.Env == "development"
}
