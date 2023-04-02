package main

import (
	"log"
	"os"
	"serverless/cli/platform"
)

func main() {
	app := platform.CreateApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
