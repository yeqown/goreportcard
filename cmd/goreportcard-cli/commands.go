package main

import (
	"os"
	"path/filepath"

	"github.com/yeqown/goreportcard/internal/types"

	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
	"github.com/yeqown/log"
)

func mountCommands(app *cli.App) {
	app.Commands = []*cli.Command{
		getStartServerCommand(),
		getCliCheckCommand(),
		getManageDBCommand(),
	}
}

func getStartServerCommand() *cli.Command {
	home, _ := os.UserHomeDir()
	confPath := filepath.Join(home, "goreportcard.toml")

	return &cli.Command{
		Name:  "start-web",
		Usage: "running web server on spec port, default=8000",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "conf",
				Usage:       "specify a path to config, default is ~/goreportcard.toml",
				Value:       confPath,
				Destination: &confPath,
			},
		},
		Action: func(c *cli.Context) error {
			log.Infof("load config file from: %s", confPath)
			if err := types.Init(confPath); err != nil {
				return errors.Wrap(err, "LoadConfig failed")
			}

			// this will blocked, it will return only if program caught an error
			return startWebServer(types.GetConfig())
		},
	}
}

func getCliCheckCommand() *cli.Command {
	var (
		dir     string
		verbose bool
	)

	return &cli.Command{
		Name:  "run",
		Usage: "running goreportcard-cli to lint project in terminal",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dir",
				Usage:       "specify an dir of golang project to run",
				Value:       ".",
				Destination: &dir,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "to show more detail about lint result",
				Destination: &verbose,
			},
		},
		Action: func(c *cli.Context) error {
			return runCli(dir, verbose)
		},
	}
}

// getManageDBCommand
// current support 2 KV DB (redis, badger),
// this command is to help user migrate data from one to another
// TODO: finish this work
func getManageDBCommand() *cli.Command {
	return &cli.Command{}
}
