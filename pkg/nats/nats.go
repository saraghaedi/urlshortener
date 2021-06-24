package nats

import (
	"strings"

	"github.com/nats-io/nats.go"
)

// Options represents a struct for creating Nats connection configurations.
type Options struct {
	Addresses []string `mapstructure:"addresses"`
}

// Create creates a nats connection.
func Create(opts Options) (*nats.Conn, error) {
	nc, err := nats.Connect(strings.Join(opts.Addresses, ","))
	if err != nil {
		return nil, err
	}

	return nc, nil
}
