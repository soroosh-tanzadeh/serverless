package worker

import (
	"context"
	"fmt"
	"sync"
)

func Worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				continue
			}

			result := job.execute(ctx)
			results <- result
			if job.Response != nil {
				job.Response <- result
			}

		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}
