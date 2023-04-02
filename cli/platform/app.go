package platform

import (
	"github.com/urfave/cli/v2"
	"serverless/cli/commnads/buid"
	"serverless/cli/commnads/run"
)

func CreateApp() *cli.App {
	return &cli.App{
		Name:  "Serverless CLI",
		Usage: "Run, Test, Deploy Serverless Apps",
		Commands: []*cli.Command{
			run.GetRunCommand(),
			buid.GetBuildCommand(),
		},
	}
}
