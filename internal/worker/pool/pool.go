package pool

import (
	"context"
	"errors"
	"serveless/internal/worker"
	"sync"
)

type WorkerPool struct {
	maxWorker int
	jobs      chan worker.Job
	wg        *sync.WaitGroup
	isClosed  bool
	context   context.Context
}

func New(maxWorker int, ctx context.Context) *WorkerPool {
	return &WorkerPool{
		maxWorker: maxWorker,
		jobs:      make(chan worker.Job, maxWorker),
		isClosed:  false,
		context:   ctx,
	}
}

func (this *WorkerPool) Run() *sync.WaitGroup {
	this.wg = &sync.WaitGroup{}
	for i := 0; i < this.maxWorker; i++ {
		this.wg.Add(1)
		go worker.Worker(this.context, this.wg, this.jobs)
	}
	go this.gracefulShutdown(this.context)
	this.wg.Wait()
	return this.wg
}

func (this *WorkerPool) gracefulShutdown(ctx context.Context) {
	go func() {
		this.wg.Add(1)
		defer this.wg.Done()
		<-ctx.Done()
		this.isClosed = true
		close(this.jobs)
	}()
}

func (this *WorkerPool) Add(job worker.Job) error {
	if !this.isClosed {
		this.jobs <- job
		return nil
	} else {
		return errors.New("adding Job to shutdown pool")
	}
}
