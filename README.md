# go-taskrunner

`go-taskrunner` is a small library for running multiple functions concurrently while
reading errors to a single `error` channel.

This library provides only a way to receive consolidated errors. It is up to the developer
to implement error handling logic, forever looping, etc.

## Installation

Install using `go get github.com/bcatubig/go-taskrunner`

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/bcatubig/go-taskrunner"
)

func main() {
    ctx := context.Background()

    t1 := func(ctx context.Context) error {
        for {
            fmt.Println("doing work")
            time.Sleep(1 * time.Second)
            fmt.Println("work done")

            select {
            case <-ctx.Done():
                return nil
            default:
            }
        }
    }

    t2 := func(ctx context.Context) error {
        // do some work
        time.Sleep(30 * time.Second)
        return nil
    )

    errs := task.RunTasks(ctx, t1, t2)

    for err := range errs {
        log.Println(err)
    }
}
```

## Testing

```shell
# run unit tests
make test

# run tests with html coverage
make test/cover
```
