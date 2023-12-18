package task

import (
	"context"
	"sync"
)

type TaskFunc func(context.Context) error

type Task struct {
	taskFn TaskFunc
	errors <-chan error
}

// NewTask creates a new task object
func NewTask(fn TaskFunc) *Task {
	t := &Task{
		taskFn: fn,
	}

	return t
}

// Run will launch the task and return a read-only error channel
func (t *Task) Run(ctx context.Context) <-chan error {
	var wg sync.WaitGroup
	wg.Add(1)

	chanError := make(chan error)

	go func() {
		defer wg.Done()

		err := t.taskFn(ctx)
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

// RunTasks will launch N number of tasks and returned a single merged error channel
// for all of the running tasks
func RunTasks(ctx context.Context, tasks ...*Task) <-chan error {
	var wg sync.WaitGroup
	errors := make(chan error)

	for _, task := range tasks {
		wg.Add(1)

		errCh := task.Run(ctx)

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
