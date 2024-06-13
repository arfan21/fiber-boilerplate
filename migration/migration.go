package migration

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var embedMigrations embed.FS

type Migration struct {
	db *sql.DB
	*migrate.Migrate
}

func New(db *sql.DB) (*Migration, error) {
	if db == nil {
		return &Migration{}, errors.New("db is nil")
	}

	source, err := iofs.New(embedMigrations, ".")
	if err != nil {
		return &Migration{}, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return &Migration{}, err
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return &Migration{}, err
	}

	return &Migration{db: db, Migrate: migrator}, nil
}
