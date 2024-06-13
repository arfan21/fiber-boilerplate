package migration

import (
	"github.com/arfan21/fiber-boilerplate/migration"
	dbpostgres "github.com/arfan21/fiber-boilerplate/pkg/db/postgres"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/urfave/cli/v2"
)

func initMigration() (*migration.Migration, error) {
	db, err := dbpostgres.NewPgx()
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(db)

	migration, err := migration.New(sqlDB)
	if err != nil {
		return nil, err
	}

	return migration, nil

}

func Root() *cli.Command {

	return &cli.Command{
		Name:  "migrate",
		Usage: "Run migration",
		Subcommands: []*cli.Command{
			Up(),
			Down(),
			Fresh(),
		},
	}
}
