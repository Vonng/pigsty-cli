package exec

import (
	"fmt"
	"testing"
)

func TestNewExecutor(t *testing.T) {
	e := NewExecutor(`/Users/vonng/pigsty/pigsty.yml`)
	fmt.Println(e)
}
