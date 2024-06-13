package main

import (
	"os"

	"github.com/arfan21/fiber-boilerplate/cmd/api"
	migration "github.com/arfan21/fiber-boilerplate/cmd/migrate"
	"github.com/arfan21/fiber-boilerplate/config"
	"github.com/urfave/cli/v2"
)

// @title fiber-boilerplate
// @version 1.0
// @description This is a sample server cell for fiber-boilerplate.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.synapsis.id
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	appCli := cli.NewApp()
	appCli.Name = config.Get().Service.Name
	appCli.Commands = []*cli.Command{
		migration.Root(),
		api.Serve(),
	}

	if err := appCli.Run(os.Args); err != nil {
		panic(err)
	}
}
