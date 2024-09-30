package config

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	DevelopmentEnvironment = "development"
	ProductionEnvironment  = "production"
)

func ReadEnv(prefix string, cfg interface{}) error {
	err := envconfig.Process(prefix, cfg)
	if err != nil {
		return err
	}
	return err
}
