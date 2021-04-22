package main

import (
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/migrations/bindata/postgres"
	"github.com/saraghaedi/urlshortener/pkg/log"
	"github.com/sirupsen/logrus"

	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/router"
	"github.com/saraghaedi/urlshortener/pkg/database"
)

func main() {
	cfg := config.Init()
	source := bindata.Resource(postgres.AssetNames(), postgres.Asset)

	log.SetupLogger(cfg.Logger.AppLogger)

	db, err := database.New(cfg.Database.Driver, cfg.Database.MasterConnStr)
	if err != nil {
		logrus.Fatalf("faled to connect to database: %s", err.Error())
	}

	if err := database.Migrate(source, cfg.Database.MasterConnStr); err != nil {
		logrus.Fatalf("faled to run database migrations: %s", err.Error())
	}

	r := router.New(cfg, db)

	logrus.Fatal(r.Start(cfg.Server.Address))
}
