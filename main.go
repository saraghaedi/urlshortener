package main

import (
	"log"
	"net/http"

	"github.com/saraghaedi/urlshortener/database"
	"github.com/saraghaedi/urlshortener/router"
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

	log.Fatal(http.ListenAndServe(":8081", r))
}
