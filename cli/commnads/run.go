package commnads

import (
	context2 "context"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"serverless/cli/server"
	"serverless/internal/engine"
	"serverless/internal/executor"
	"serverless/internal/packer"
	"serverless/internal/utils/file"
	"syscall"
	"time"
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
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   8090,
				Usage:   "Webserver Port",
			},
			&cli.StringFlag{
				Name:  "host",
				Value: "127.0.0.1",
				Usage: "Webserver Host",
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
			httpServer := server.NewInternalHttpServer(context.String("host"), context.Int("port"), folder+"/build", manifest)
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
			go func() {
				err := httpServer.Start()
				if err != nil && err != http.ErrServerClosed {
					log.Fatal(err.Error())
				}
			}()

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			<-sigs

			log.Println("Graceful shutdown within 5sec...")
			timeoutContext, cancel := context2.WithTimeout(context.Context, 5*time.Second)
			err = httpServer.Stop(timeoutContext)
			if err != nil {
				log.Fatal(err)
			}

			defer cancel()
			<-timeoutContext.Done()

			return nil
		},
	}
}
