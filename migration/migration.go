package migration

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

type Migration struct {
	db *sql.DB
}

func New(db *sql.DB) (*Migration, error) {
	if db == nil {
		return &Migration{}, errors.New("db is nil")
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return &Migration{}, err
	}

	return &Migration{db: db}, nil
}

func (m *Migration) Up(ctx context.Context) error {
	return goose.UpContext(ctx, m.db, ".")
}

func (m *Migration) Down(ctx context.Context) error {
	return goose.ResetContext(ctx, m.db, ".")
}
