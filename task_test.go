package task_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	task "github.com/bcatubig/go-taskrunner"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func testTask(sleepDuration time.Duration) task.TaskFunc {
	taskID := rand.Intn(100)
	logger := hclog.New(&hclog.LoggerOptions{
		Name: fmt.Sprintf("task-%d", taskID),
	})

	return func(ctx context.Context) error {
		for {
			logger.Info("doing work")
			luckyNum := rand.Intn(100)
			if luckyNum%2 == 0 {
				return fmt.Errorf("task-%d: UH OH: something went wrong", taskID)
			}
			select {
			case <-ctx.Done():
				logger.Info("context cancelled")

				return nil
			default:
				time.Sleep(sleepDuration)
			}
		}
	}
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	t.Run("happy path", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		t1 := task.NewTask(func(ctx context.Context) error { return nil })
		errs := task.RunTasks(ctx, t1)

		for err := range errs {
			assert.Nil(t, err)
		}
	})

	t.Run("task error", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		t1 := task.NewTask(func(ctx context.Context) error { return fmt.Errorf("uh oh!") })
		errs := task.RunTasks(ctx, t1)

		for err := range errs {
			assert.NotNil(t, err)
		}
	})

	t.Run("forever task with cancellation", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		ctx, cancel := context.WithTimeout(ctx, 4*time.Millisecond)
		defer cancel()

		t1 := task.NewTask(func(ctx context.Context) error {
			for {
				time.Sleep(1 * time.Millisecond)
				select {
				case <-ctx.Done():
					return nil
				default:
				}
			}
		})
		errs := task.RunTasks(ctx, t1)

		for err := range errs {
			assert.Nil(t, err)
		}
	})
}
