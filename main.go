package main

import (
	"github.com/saraghaedi/urlshortener/database"
	"github.com/saraghaedi/urlshortener/router"
	"log"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatalf("faled to connect to database: %s", err.Error())
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("faled to run database migrations: %s", err.Error())
	}

	r := router.New(db)

	log.Fatal(r.Start(":8080"))
}
