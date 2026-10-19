package collector

import (
	"context"
	"large-file-processing/model"
)

type ErrorCollector struct {
	ctx context.Context
	Errors chan model.ProcessingError
}

func NewErrorCollector(ctx context.Context, bufferSize int) *ErrorCollector {
	return &ErrorCollector{
		ctx: ctx,
		Errors: make(chan model.ProcessingError, bufferSize),
	}
}

func (ec *ErrorCollector) Add(err model.ProcessingError) error {
	select {
	case <-ec.ctx.Done():
		return ec.ctx.Err()
	case ec.Errors <- err:
		return nil
	}
}