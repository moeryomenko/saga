package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// LoadConfig returns parsed from environment variables service configuration.
func LoadConfig() (*Config, error) {
	conf := &Config{}
	err := envconfig.Process("", conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// Config represents service configurations.
type Config struct {
	Health HealthConfig `envconfig:"HEALTH"`
}

// PoolConfig represents databese connection pool configuration.
type PoolConfig struct {
	MaxOpenConns int `envconfig:"MAX_OPEN_CONNS" default:"20"`
	MaxIdleConns int `envconfig:"MAX_IDLE_CONNS" default:"20"`
}

// HealthConfig represents health controller configuration.
type HealthConfig struct {
	Port          int           `envconfig:"PORT" default:"6061"`
	LiveEndpoint  string        `envconfig:"LIVINESS_ENDPOINT" default:"/livez"`
	ReadyEndpoint string        `envconfig:"READINESS_ENDPOINT" default:"/ready"`
	Period        time.Duration `envconfig:"PERIOD" default:"3s"`
	GracePeriod   time.Duration `envconfig:"GRACE_PERIOD" default:"30s"`
}
