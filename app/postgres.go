package app

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
)

func getDBConnection() (*sql.DB, error) {
	connectionString := getConnectionString()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		logrus.WithError(err).Fatal("Error: The data source args are not valid")
	}

	err = db.Ping()
	if err != nil {
		logrus.WithField("error", err).Fatal("error: Could not establish a connection with the database")
	}

	return db, err
}

func getConnectionString() string {
	return "port=5432 host=127.0.0.1 dbname=postgres user=postgres password=Password123! sslmode=disable"
}
