package config

import (
	"time"

	"github.com/saraghaedi/urlshortener/pkg/redis"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sirupsen/logrus"

	"github.com/saraghaedi/urlshortener/pkg/prometheus"

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
		Logger     Logger     `mapstructure:"logger"`
		Server     Server     `mapstructure:"server"`
		Redis      Redis      `mapstructure:"redis"`
		Database   Database   `mapstructure:"database"`
		Monitoring Monitoring `mapstructure:"monitoring"`
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

	// Redis represents Redis configuration struct.
	Redis struct {
		MasterAddress string        `mapstructure:"master-address"`
		SlaveAddress  string        `mapstructure:"slave-address"`
		Options       redis.Options `mapstructure:"options"`
	}

	// Monitoring represents monitoring configuration struct.
	Monitoring struct {
		Prometheus prometheus.Config `mapstructure:"prometheus"`
	}
)

// Validate validates Database struct.
func (d Database) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(
			&d.Driver,
			validation.In("postgres"),
		),
	)
}

// Validate validates Config struct.
func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Database,
		),
	)
}

// Init initializes application configuration.
func Init() Config {
	var cfg Config

	config.Init(app, cfgFile, &cfg, defaultConfig, cfgPrefix)

	if err := cfg.Validate(); err != nil {
		logrus.Fatalf("failed to validate configurations: %s", err.Error())
	}

	return cfg
}
