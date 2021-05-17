package server

import (
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/router"
	"github.com/saraghaedi/urlshortener/pkg/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	masterDb, err := database.New(cfg.Database.Driver, cfg.Database.MasterConnStr)
	if err != nil {
		logrus.Fatalf("faled to connect to master database: %s", err.Error())
	}

	slaveDb, err := database.New(cfg.Database.Driver, cfg.Database.SlaveConnStr)
	if err != nil {
		logrus.Fatalf("faled to connect to slave database: %s", err.Error())
	}

	r := router.New(cfg, masterDb, slaveDb)
	logrus.Fatal(r.Start(cfg.Server.Address))
}

// Register registers server command for urlshortener binary.
func Register(root *cobra.Command, cfg config.Config) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run url shortener server component",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg)
		},
	}
	root.AddCommand(cmd)
}
