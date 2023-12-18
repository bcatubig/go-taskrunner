package task

import (
	"context"
	"sync"
)

type TaskFunc func(context.Context) error

// RunTasks will launch N number of tasks and returned a single merged error channel
// for all of the running tasks
func RunTasks(ctx context.Context, tasks ...TaskFunc) <-chan error {
	var wg sync.WaitGroup
	errors := make(chan error)

	for _, task := range tasks {
		wg.Add(1)

		errCh := runTask(ctx, task)

		go func(c <-chan error) {
			defer wg.Done()
			for err := range c {
				errors <- err
			}
		}(errCh)

	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	return errors
}

// runTasks will launch the task and return a read-only error channel
func runTask(ctx context.Context, fn TaskFunc) <-chan error {
	var wg sync.WaitGroup
	wg.Add(1)

	chanError := make(chan error)

	go func() {
		defer wg.Done()

		err := fn(ctx)
		if err != nil {
			chanError <- err
		}
	}()

	go func() {
		wg.Wait()
		close(chanError)
	}()

	return chanError
}
