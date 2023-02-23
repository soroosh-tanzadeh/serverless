package commnads

import (
	"errors"
	"github.com/urfave/cli/v2"
	"serveless/internal/packer"
	"serveless/internal/utils/file"
)

func GetBuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build Serverless Application",
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
			return nil
		},
	}
}
