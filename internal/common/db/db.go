package db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/morzhanov/go-realworld/internal/common/config"
)

func NewDb(c *config.Config) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", c.PsqlConnectionString)
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
	return m.Steps(2)
}
