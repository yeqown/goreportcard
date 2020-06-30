package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Authors = []*cli.Author{
		{
			Name:  "yeqown",
			Email: "yeqown@gmail.com",
		},
	}
	app.Copyright = "2020@yeqown"

	mountCommands(app)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
