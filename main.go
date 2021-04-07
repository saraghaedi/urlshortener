package main

import (
	"github.com/saraghaedi/urlshortener/log"
	"github.com/sirupsen/logrus"

	"github.com/saraghaedi/urlshortener/database"
	"github.com/saraghaedi/urlshortener/router"
)

func main() {
	log.SetupLogger(log.AppLogger{
		Level:  logrus.InfoLevel.String(),
		StdOut: true,
	})

	db, err := database.New()
	if err != nil {
		logrus.Fatalf("faled to connect to database: %s", err.Error())
	}

	if err := database.Migrate(db); err != nil {
		logrus.Fatalf("faled to run database migrations: %s", err.Error())
	}

	r := router.New(db)

	logrus.Fatal(r.Start(":8080"))
}
