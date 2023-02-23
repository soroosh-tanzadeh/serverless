package pool

import (
	"context"
	"serveless/internal/worker"
	"sync"
)

type WorkerPool struct {
	maxWorker int
	jobs      chan worker.Job
	results   chan worker.Result
	Done      chan struct{}
}

func New(maxWorker int) (*WorkerPool, error) {
	return &WorkerPool{
		maxWorker: maxWorker,
		jobs:      make(chan worker.Job, maxWorker),
		results:   make(chan worker.Result, maxWorker),
		Done:      make(chan struct{}),
	}, nil
}

func (this *WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < this.maxWorker; i++ {
		wg.Add(1)
		go worker.Worker(ctx, &wg, this.jobs, this.results)
	}
	wg.Wait()
	close(this.Done)
	close(this.results)
}

func (this WorkerPool) Results() <-chan worker.Result {
	return this.results
}

func (this WorkerPool) Add(job worker.Job) {
	this.jobs <- job
	close(this.jobs)
}
