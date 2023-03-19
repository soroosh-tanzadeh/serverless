package pool

import (
	"context"
	"errors"
	"serverless/internal/worker"
	"sync"
)

type WorkerPool struct {
	maxWorker int
	jobs      chan worker.Job
	wg        *sync.WaitGroup
	isClosed  bool
}

func New(maxWorker int) *WorkerPool {
	return &WorkerPool{
		maxWorker: maxWorker,
		jobs:      make(chan worker.Job, maxWorker),
		isClosed:  false,
	}
}

func (w *WorkerPool) Run(context context.Context) *sync.WaitGroup {
	w.wg = &sync.WaitGroup{}
	for i := 0; i < w.maxWorker; i++ {
		w.wg.Add(1)
		go worker.Worker(context, w.wg, w.jobs)
	}
	w.gracefulShutdown(context)
	w.wg.Wait()
	return w.wg
}

func (w *WorkerPool) gracefulShutdown(ctx context.Context) {
	go func() {
		w.wg.Add(1)
		defer w.wg.Done()
		<-ctx.Done()
		w.isClosed = true
		close(w.jobs)
	}()
}

func (w *WorkerPool) Add(job worker.Job) error {
	if w.isClosed {
		return errors.New("adding job to stopped worker")
	}
	select {
	case w.jobs <- job:
		return nil
	default:
		return errors.New("error adding job to channel")
	}
}
