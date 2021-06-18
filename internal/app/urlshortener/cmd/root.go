package cmd

import (
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/cmd/migrate"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/cmd/server"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/cmd/sync"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/pkg/log"
	"github.com/spf13/cobra"
)

// NewRootCommand creates a new urlshortener root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:   "urlshortener",
		Short: "URLShortener is an app for shortening URLs",
	}

	cfg := config.Init()

	log.SetupLogger(cfg.Logger.AppLogger)

	migrate.Register(root, cfg)
	server.Register(root, cfg)
	sync.Register(root, cfg)

	return root
}
