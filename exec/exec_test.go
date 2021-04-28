package exec

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewExecutor(t *testing.T) {
	e := NewExecutor(`/Users/vonng/pigsty/pigsty.yml`)
	fmt.Println(e)
}

func TestExecutor_NewJob(t *testing.T) {
	e := NewExecutor(`/Users/vonng/pigsty/pigsty.yml`)
	job := e.NewJob(
		WithPlaybook("pgsql.yml"),
		WithName("pgsql init"),
		WithLogPath("/tmp/test.log"),
	)

	print(job.JSON())
	go func() {
		time.Sleep(time.Second * 2)
		job.Cancel()
	}()
	job.Run(context.TODO())
}
