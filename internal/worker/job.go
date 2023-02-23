package worker

import (
	"context"
	"github.com/google/uuid"
)

type JobID string
type ExecutionFn func(ctx context.Context, args interface{}) (interface{}, error)

type JobDescriptor struct {
	ID JobID
}

type Job struct {
	Args       interface{}
	Action     ExecutionFn
	Descriptor JobDescriptor
	Response   chan Result
}

type Result struct {
	Value      interface{}
	Err        error
	Descriptor JobDescriptor
}

func New(args interface{}, action ExecutionFn, response chan Result) *Job {
	return &Job{
		Descriptor: JobDescriptor{
			ID: JobID(uuid.NewString()),
		},
		Args:     args,
		Action:   action,
		Response: response,
	}
}

func (this *Job) execute(ctx context.Context) Result {
	value, err := this.Action(ctx, this.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: this.Descriptor,
		}
	}

	return Result{
		Value:      value,
		Descriptor: this.Descriptor,
	}
}
