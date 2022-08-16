package config

import (
	"fmt"
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
	Host string `envconfig:"HOST"`
	Port int    `envconfig:"PORT" default:"8080"`

	Health HealthConfig `envconfig:"HEALTH"`

	Database *DBConfig `envconfig:"DB"`
}

// Addr returns address for listening.
func (cfg *Config) Addr() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

// SessionConfig represents configuration for auth middleware.
type SessionConfig struct {
	Secret         string        `envconfig:"SECRET"`
	ExpirationTime time.Duration `envconfig:"EXPIRATION" default:"5m"`
}

// DBConfig represents database connection configuration.
type DBConfig struct {
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT" default:"5432"`
	Name     string `envconfig:"NAME"`
	User     string `envconfig:"USER"`
	Password string `encconfig:"PASSWORD"`

	Pool *PoolConfig `envconfig:"POOL"`
}

// PoolConfig represents databese connection pool configuration.
type PoolConfig struct {
	MaxOpenConns int `envconfig:"MAX_OPEN_CONNS" default:"20"`
	MaxIdleConns int `envconfig:"MAX_IDLE_CONNS" default:"20"`
}

// HealthConfig represents health controller configuration.
type HealthConfig struct {
	Port          int           `envconfig:"PORT" default:"6060"`
	LiveEndpoint  string        `envconfig:"LIVINESS_ENDPOINT" default:"/livez"`
	ReadyEndpoint string        `envconfig:"READINESS_ENDPOINT" default:"/ready"`
	Period        time.Duration `envconfig:"PERIOD" default:"3s"`
	GracePeriod   time.Duration `envconfig:"GRACE_PERIOD" default:"30s"`
}
