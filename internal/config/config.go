package config

import (
	"os"

	"github.com/akash202k/pulse/internal/model"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Services []model.Service `yaml:"services"`
	SLO      model.SLO       `yaml:"SLO"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	// if err := yaml.Unmarshal(data, &cfg); err != nil {
	// 	return nil, err
	// }
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
