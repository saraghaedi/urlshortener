package main

import (
	"github.com/saraghaedi/urlshortener/config"
	"github.com/saraghaedi/urlshortener/pkg/log"
	"github.com/sirupsen/logrus"

	"github.com/saraghaedi/urlshortener/database"
	"github.com/saraghaedi/urlshortener/router"
)

func main() {
	cfg := config.Init()

	log.SetupLogger(cfg.Logger.AppLogger)

	db, err := database.New(cfg.Database.Driver, cfg.Database.MasterConnStr)
	if err != nil {
		logrus.Fatalf("faled to connect to database: %s", err.Error())
	}

	if err := database.Migrate(db); err != nil {
		logrus.Fatalf("faled to run database migrations: %s", err.Error())
	}

	r := router.New(cfg, db)

	logrus.Fatal(r.Start(cfg.Server.Address))
}
