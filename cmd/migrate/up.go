package migration

import (
	"github.com/urfave/cli/v2"
)

func Up() *cli.Command {

	return &cli.Command{
		Name:  "up",
		Usage: "Run all migrations",
		Action: func(c *cli.Context) error {
			migrations, err := initMigration()
			if err != nil {
				return err
			}

			if err := migrations.Up(c.Context); err != nil {
				return err
			}

			return nil
		},
	}
}
