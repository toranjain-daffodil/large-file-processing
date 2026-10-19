package pipeline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"large-file-processing/collector"
	"large-file-processing/model"
	"large-file-processing/reader"
	"large-file-processing/worker"
	"os"
	"sync"
)

const (
	workerCount = 8
	jobsBuffer = 1000
	errorBuffer = 256
)

func Run(
	parent context.Context,
	inputPath string,
	errorPath string,
	stats *collector.StatisticsCollector,
) error {
	if stats == nil {
		return errors.New("statistics collector is nil")
	}

	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	jobs := make(chan model.RawEmployee, jobsBuffer)
	errorCollector := collector.NewErrorCollector(ctx, errorBuffer)

	fatalErrors := make(chan error, workerCount+2)

	reportFatal := func(err error) {
		if err == nil || errors.Is(err, context.Canceled) {
			return
		}

		cancel()
		fatalErrors <- err
	}

	writerDone := make(chan struct{})

	go func(){
		defer close(writerDone)

		file, err := os.Create(errorPath)
		if err != nil {
			reportFatal(fmt.Errorf("create error report: %w", err))
			return
		}

		encoder := json.NewEncoder(file)

		for record := range errorCollector.Errors {
			if err := encoder.Encode(record); err != nil {
				_ = file.Close()
				reportFatal(fmt.Errorf("write error report: %w", err))
				return
			}
		}

		if err := file.Close(); err != nil {
			reportFatal(fmt.Errorf("close error report: %w", err))
		}
	}()

	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i:=0; i < workerCount; i++ {
		go func(){
			defer wg.Done()

			for {
				if ctx.Err() != nil {
					return
				}

				select {
				case <-ctx.Done():
					return
				case raw, ok := <-jobs:
					if !ok {
						return
					}
					if err := worker.ProcessEmployee(ctx, raw, stats, errorCollector); err != nil {
						reportFatal(err)
						return
					}
				}
			}
		}()
	}

	producerDone := make(chan error, 1)

	go func() {
		defer close(jobs)

		err := reader.ReadCSV(inputPath, func(raw model.RawEmployee) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case jobs<-raw:
				return nil
			}
		})

		if err != nil && ctx.Err() == nil {
			reportFatal(fmt.Errorf("read csv: %w", err))
		}

		producerDone <- err
	}()

	var producerErr error

	select {
	case producerErr = <- producerDone:
	case <-ctx.Done():
		producerErr = <-producerDone
	}

	wg.Wait()

	close(errorCollector.Errors)
	<-writerDone

	var fatalErr error
	
	select {
	case fatalErr = <-fatalErrors:
	default:
	}

	if fatalErr != nil {
		return fmt.Errorf("pipeline failed; results are incomplete: %w", fatalErr)
	}

	if err := parent.Err(); err != nil {
		return fmt.Errorf("pipeline cancelled: %w", err)
	}
	if producerErr != nil {
		return fmt.Errorf("csv producer failed: %w", producerErr)
	}

	return nil
}