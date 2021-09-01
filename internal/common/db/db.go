package db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func NewDb() (*sqlx.DB, error) {
	psqlConnectionString := viper.GetString("PSQL_CONNECTION_STRING")
	return sqlx.Connect("postgres", psqlConnectionString)
}

func RunMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	m.Steps(2)
	return nil
}
