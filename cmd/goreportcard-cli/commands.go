package main

import (
	"github.com/urfave/cli"
)

func mountCommands(app *cli.App) {
	app.Commands = []cli.Command{
		getStartSevrerCommand(),
		getCliCheckCommand(),
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
	var (
		dir     string
		verbose bool
	)

	return cli.Command{
		Name:  "run",
		Usage: "running goreportcard-cli to lint project in terminal",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "dir",
				Usage:       "specify an dir of golang project to run",
				Value:       ".",
				Destination: &dir,
			},
			cli.BoolFlag{
				Name:        "verbose",
				Usage:       "to show more detail about lint result",
				Destination: &verbose,
			},
		},
		Action: func(c *cli.Context) error {
			return cliCheck(dir, verbose)
		},
	}
}

// TODO:
func getManageDBCommand() cli.Command {
	return cli.Command{}
}
