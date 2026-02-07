package config

import "github.com/akash202k/pulse/internal/model"

type Config struct {
	Services []model.Service `yaml:"services"`
	SLO      model.SLO       `yaml:"SLO"`
}
