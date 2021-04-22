package config

import (
	"time"

	"github.com/saraghaedi/urlshortener/pkg/config"
	"github.com/saraghaedi/urlshortener/pkg/log"
)

const (
	app       = "urlshortener"
	cfgFile   = "config.yaml"
	cfgPrefix = "urlshortener"
)

type (
	// Config represents application configuration struct.
	Config struct {
		Logger   Logger   `mapstructure:"logger"`
		Server   Server   `mapstructure:"server"`
		Database Database `mapstructure:"database"`
	}

	// Logger represents logger configuration struct.
	Logger struct {
		AccessLogger log.AccessLogger `mapstructure:"access"`
		AppLogger    log.AppLogger    `mapstructure:"app"`
	}

	// Server represents server configuration struct.
	Server struct {
		Address         string        `mapstructure:"address"`
		ReadTimeout     time.Duration `mapstructure:"read-timeout"`
		WriteTimeout    time.Duration `mapstructure:"write-timeout"`
		GracefulTimeout time.Duration `mapstructure:"graceful-timeout"`
	}

	// Database represents database configuration struct.
	Database struct {
		Driver        string `mapstructure:"driver"`
		MasterConnStr string `mapstructure:"master-conn-string"`
		SlaveConnStr  string `mapstructure:"slave-conn-string"`
	}
)

// Init initializes application configuration.
func Init() Config {
	var cfg Config

	config.Init(app, cfgFile, &cfg, defaultConfig, cfgPrefix)

	return cfg
}
