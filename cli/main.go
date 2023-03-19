package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"serverless/cli/commnads"
)

func main() {
	app := &cli.App{
		Name:  "Serverless CLI",
		Usage: "Run, Test, Deploy Serverless Apps",
		Commands: []*cli.Command{
			commnads.GetRunCommand(),
			commnads.GetBuildCommand(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
