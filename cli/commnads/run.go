package commnads

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
	"log"
	"serveless/cli/server"
	"serveless/internal/engine"
	"serveless/internal/executor"
	"serveless/internal/packer"
	"serveless/internal/utils/file"
)

func build(folder string) (*engine.Manifest, error) {
	_, err := packer.Build(folder)
	if err != nil {
		log.Fatal(err)
	}
	manifest, err := executor.ParseManifest(folder)
	if err != nil {
		log.Fatal(err)
	}
	return manifest, err
}

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
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatal(err)
			}
			defer watcher.Close()

			var folder string
			if context.NArg() > 0 {
				folder = context.Args().Get(0)
			} else {
				return errors.New("application path is required")
			}
			if isDir, _ := file.IsDirectory(folder); !isDir {
				return errors.New("selected path is not folder")
			}

			manifest, err := build(folder)
			if err != nil {
				return err
			}
			httpServer := server.NewInternalHttpServer("0.0.0.0", 8090, folder+"/build", manifest)
			log.Println("Running Http Server on Port 8090...")
			// Start listening for events.
			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}
						if event.Has(fsnotify.Write) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
							_, err := build(folder)
							if err != nil {
								log.Println(err.Error())
							}
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}
						log.Println("error:", err)
					}
				}
			}()
			err = watcher.Add(folder)
			if err != nil {
				log.Fatal(err)
			}
			err = httpServer.Start()
			if err != nil {
				return err
			}
			<-make(chan struct{})
			return nil
		},
	}
}
