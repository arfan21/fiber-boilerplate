package migration

import (
	"github.com/urfave/cli/v2"
)

func Down() *cli.Command {

	return &cli.Command{
		Name:  "down",
		Usage: "Rollback the last migration",
		Action: func(c *cli.Context) error {
			migrations, err := initMigration()
			if err != nil {
				return err
			}

			if err := migrations.Down(c.Context); err != nil {
				return err
			}

			return nil
		},
	}
}
