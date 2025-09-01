// Package slicex provides utility functions for slice operations with generics support.
package slicex

import (
	"context"
	"errors"
	"sync"
)

// Unique returns a new slice containing only unique elements from the input slice,
// preserving the order of first occurrence.
func Unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}

	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// FilterNonZero returns a new slice with all non-zero values from the input slice.
// Zero values are determined by Go's zero value concept (0, "", nil, etc.).
func FilterNonZero[T comparable](slice []T) []T {
	var zero T
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if item != zero {
			result = append(result, item)
		}
	}

	return result
}

// Map applies the given function to each element of the slice and returns
// a new slice containing the results.
func Map[T, R any](slice []T, fn func(T) R) []R {
	if len(slice) == 0 {
		return nil
	}

	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}

	return result
}

// Group groups the elements of the slice by the mapConcurrentResult of the key function.
// Returns a map where keys are the grouping criteria and values are slices
// of grouped items.
func Group[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)

	for _, item := range slice {
		key := keyFn(item)
		result[key] = append(result[key], item)
	}

	return result
}

// MapConcurrentHandler provides fluent configuration for concurrent map operations.
type MapConcurrentHandler[T, R any] struct {
	mapFunc     func(context.Context, T) (R, error)
	concurrency int
	stopOnError bool
}

// WithConcurrency sets the maximum number of concurrent operations.
// Defaults to 8 if not specified.
func (h *MapConcurrentHandler[T, R]) WithConcurrency(n int) *MapConcurrentHandler[T, R] {
	h.concurrency = n
	return h
}

// WithStopOnError configures whether to stop processing on first error (true)
// or collect all errors and continue processing (false).
// Defaults to true (stop on first error).
func (h *MapConcurrentHandler[T, R]) WithStopOnError(stop bool) *MapConcurrentHandler[T, R] {
	h.stopOnError = stop
	return h
}

// mapConcurrentJob represents a work item for the worker pool
type mapConcurrentJob[T any] struct {
	index int
	value T
}

// mapConcurrentResult represents the mapConcurrentResult of processing a mapConcurrentJob
type mapConcurrentResult[R any] struct {
	index int
	value R
	err   error
}

// Execute runs the concurrent map operation on the provided slice.
// Returns a slice of results preserving input order and any errors encountered.
func (h *MapConcurrentHandler[T, R]) Execute(ctx context.Context, items []T) ([]R, error) {
	if len(items) == 0 {
		return nil, nil
	}

	// Determine actual number of workers (min of concurrency and items length)
	numWorkers := h.concurrency
	if n := len(items); n < numWorkers {
		numWorkers = n
	}

	// Pre-allocate mapConcurrentResult items to preserve ordering
	results := make([]R, len(items))
	errs := make([]error, len(items)+1)

	// Create channels for mapConcurrentJob distribution and mapConcurrentResult collection
	jobs := make(chan mapConcurrentJob[T], len(items))

	// Context for cancellation on first error
	child, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	startWorker := func() {
		defer wg.Done()
		for {
			select {
			case <-child.Done():
				return

			case item, ok := <-jobs:
				if !ok {
					return
				}
				v, err := h.mapFunc(ctx, item.value)
				if err != nil {
					errs[item.index] = err
					if h.stopOnError {
						cancel()
						return
					}
				} else {
					results[item.index] = v
				}
			}
		}
	}

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go startWorker()
	}

	// Send all jobs to workers
	go func() {
		defer close(jobs)
		for i, item := range items {
			select {
			case jobs <- mapConcurrentJob[T]{index: i, value: item}:
			case <-child.Done():
				return
			}
		}
	}()

	// wait for all workers to complete
	wg.Wait()
	errs = append(errs, ctx.Err()) // ctx.Err is nil if no error
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return results, nil
}

// MapConcurrent creates a new concurrent map handler with the given mapping function.
// The mapping function should have the signature: func(context.Context, T) (R, error).
// Returns a handler that can be configured with fluent methods before execution.
func MapConcurrent[T, R any](mapFunc func(context.Context, T) (R, error)) *MapConcurrentHandler[T, R] {
	return &MapConcurrentHandler[T, R]{
		mapFunc:     mapFunc,
		concurrency: 8,    // Default concurrency level
		stopOnError: true, // Default behavior: stop on first error
	}
}
