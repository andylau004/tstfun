package syncx

import (
	"fmt"
	"testing"
	// "timex"
)

func TestImmutableResource(t *testing.T) {
	var count int
	r := NewImmutableResource(func() (interface{}, error) {
		//r := timex.NewImmutableResource(func() (interface{}, error) {
		count++
		return "hello", nil
	})

	res, err := r.Get()
	fmt.Println("res=", res, ", err=", err)

	{
		res, err := r.Get()
		fmt.Println("res=", res, ", err=", err)

	}
}
