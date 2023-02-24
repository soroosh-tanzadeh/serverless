package worker

import (
	"context"
	"sync"
)

func Worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job) {
	defer wg.Done()
	for job := range jobs {
		result := job.execute(ctx)
		if job.Response != nil {
			job.Response <- result
		}
	}
}
