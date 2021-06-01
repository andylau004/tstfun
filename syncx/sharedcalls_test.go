package syncx

import (
	"fmt"
	"testing"
)

func TestExclusiveCallDo(t *testing.T) {
	g := NewSharedCalls()
	v, err := g.Do("key", func() (interface{}, error) {
		return "bar", nil
	})
	_ = err
	// got := fmt.Sprintf("%v (%T)", v, v)
	// fmt.Printf("got: %s\n", got)
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got == want {
		t.Errorf("Do = %v; want = %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}

	xx := fmt.Sprintf("%T", v)
	fmt.Println("xx=", xx)
}
