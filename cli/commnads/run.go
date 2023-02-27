package commnads

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"serveless/internal/executor"
	"serveless/internal/packer"
	"serveless/internal/utils/file"
)

func GetRunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Execute Serverless Application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "folder",
				Aliases: []string{"f"},
				Usage:   "Application Folder",
			},
		},
		Action: func(context *cli.Context) error {
			var folder string
			if context.NArg() > 0 {
				folder = context.Args().Get(0)
			} else {
				return errors.New("application path is required")
			}

			if isDir, _ := file.IsDirectory(folder); !isDir {
				return errors.New("selected path is not folder")
			}
			_, err := packer.Build(folder)
			if err != nil {
				return err
			}
			manifest, err := executor.ParseManifest(folder)
			if err != nil {
				return err
			}
			responseChannel, err := executor.Execute(folder+"/build", manifest)
			if err != nil {
				return err
			}
			response := <-responseChannel
			fmt.Printf("Http Status Code: %d\n", response.Status)

			fmt.Printf("%s\n", response.Content)

			fmt.Print("Headers: ")
			fmt.Println(response.Headers)

			return nil
		},
	}
}
