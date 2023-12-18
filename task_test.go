package task_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	task "github.com/bcatubig/go-taskrunner"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx := context.Background()
	t.Run("happy path", func(t *testing.T) {
		t1 := func(ctx context.Context) error { return nil }
		errs := task.RunTasks(ctx, t1)

		for err := range errs {
			assert.Nil(t, err)
		}
	})

	t.Run("task error", func(t *testing.T) {
		t1 := func(ctx context.Context) error { return fmt.Errorf("uh oh!") }
		errs := task.RunTasks(ctx, t1)

		for err := range errs {
			assert.NotNil(t, err)
		}
	})

	t.Run("forever task with cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 4*time.Millisecond)
		defer cancel()

		t1 := func(ctx context.Context) error {
			for {
				time.Sleep(1 * time.Millisecond)
				select {
				case <-ctx.Done():
					return nil
				default:
				}
			}
		}

		errs := task.RunTasks(ctx, t1)

		for err := range errs {
			assert.Nil(t, err)
		}
	})
}
