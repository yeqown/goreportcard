package main

import "github.com/urfave/cli"

func mountCommands(app *cli.App) {
	app.Commands = []cli.Command{
		getStartSevrerCommand(),
		getCliCheckCommand(),
		getCleanRepoCommand(),
		getManageDBCommand(),
	}
}

func getStartSevrerCommand() cli.Command {
	var port int

	return cli.Command{
		Name:  "start-web",
		Usage: "running web server on spec port, default=8000",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "port",
				Usage:       "which port to listen",
				Value:       8000,
				Destination: &port,
			},
		},
		Action: func(c *cli.Context) error {
			// this will blocked, it will return only if program caught an error
			return startWebServer(c.Int("port"))
		},
	}
}

func getCliCheckCommand() cli.Command {
	return cli.Command{
		Action: func(c *cli.Context) error {
			cliCheck()
			return nil
		},
	}
}

// TODO:
func getCleanRepoCommand() cli.Command {
	return cli.Command{}
}

// TODO:
func getManageDBCommand() cli.Command {
	return cli.Command{}
}
