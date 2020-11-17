package database

import "fmt"

const (
	dbUser     = "urldb"
	dbPassword = "secret"
	dbName     = "urldb"
	dbHost     = "0.0.0.0"
	dbPort     = "5432"
	dbType     = "postgres"
)

func connectionString() string {
	dataBase := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName, dbPassword,
	)

	return dataBase
}
