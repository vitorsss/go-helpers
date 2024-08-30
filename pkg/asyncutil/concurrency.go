package asyncutil

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

func ConcurrencyExec[T any, R any](
	ctx context.Context,
	items []T,
	fn func(ctx context.Context, item T) (R, error),
	options ...ConcurrencyExecOption,
) ([]R, error) {
	opt := defaultOptions()
	for _, option := range options {
		option(opt)
	}

	resultChan := make(chan execPair[R], opt.maxConcurrency)

	results := make([]R, 0, len(items))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opened := 0
	wg := &sync.WaitGroup{}
	var lastErr error
	for _, item := range items {
		if opened >= opt.maxConcurrency {
			result := <-resultChan
			opened--
			if result.err != nil {
				if !opt.waitPartialSuccess {
					return nil, result.err
				}
				lastErr = result.err
				break
			} else {
				results = append(results, result.result)
			}
		}
		opened++
		wg.Add(1)
		go safeExec(ctx, item, resultChan, wg, fn)
	}

	wg.Wait()
	close(resultChan)
	for len(resultChan) > 0 {
		result := <-resultChan
		if result.err != nil {
			if !opt.waitPartialSuccess {
				return nil, result.err
			}
			lastErr = result.err
		} else {
			results = append(results, result.result)
		}
	}

	return results, lastErr
}

type execPair[R any] struct {
	result R
	err    error
}

func safeExec[T any, R any](
	ctx context.Context,
	item T,
	resultChan chan execPair[R],
	wg *sync.WaitGroup,
	fn func(ctx context.Context, item T) (R, error),
) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			resultChan <- execPair[R]{
				err: errors.Errorf("recovered error on safe exec: %v", err),
			}
		}
	}()
	result, err := fn(ctx, item)

	resultChan <- execPair[R]{
		result: result,
		err:    err,
	}
}

func defaultOptions() *options {
	return &options{
		maxConcurrency:     10,
		waitPartialSuccess: false,
	}
}

type options struct {
	maxConcurrency     int
	waitPartialSuccess bool
}

type ConcurrencyExecOption func(opt *options)

func WithMaxConcurrency(maxConcurrency int) ConcurrencyExecOption {
	return func(opt *options) {
		opt.maxConcurrency = maxConcurrency
	}
}

func WithWaitPartialSuccess() ConcurrencyExecOption {
	return func(opt *options) {
		opt.waitPartialSuccess = true
	}
}
