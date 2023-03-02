package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"serveless/internal/worker"
	"serveless/internal/worker/pool"
	"syscall"
	"time"
)

func main() {
	appContext, cancel := context.WithCancel(context.TODO())
	workerPool := pool.New(100)
	go app(workerPool)

	// Handel Interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		for {
			sig := <-sigs
			if sig == syscall.SIGINT {
				cancel()
			}
		}
	}()
	workerPool.Run(appContext)
}

func app(workerPool *pool.WorkerPool) {
	fmt.Println("Main App Thread")
	for i := 0; i < 20; i++ {
		var jobId = i
		time.Sleep(time.Second * 1)
		workerPool.Add(*worker.NewJob(nil, func(ctx context.Context, args interface{}) (interface{}, error) {
			fmt.Printf("Job %d Started!\n", jobId)
			time.Sleep(time.Second * 5)
			fmt.Printf("Job %d Done!\n", jobId)
			return i, nil
		}, nil))
	}
}
