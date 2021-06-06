package syncx

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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

func TestExclusiveCallDoDupSuppress(t *testing.T) {
	g := NewSharedCalls()

	c := make(chan string)
	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	var wg sync.WaitGroup
	wg.Add(10)

	const count = 10
	for i := 0; i < count; i++ {

		go func() {
			ret, err := g.Do("key", fn)
			if ret.(string) != "bar" {
				fmt.Println("fn ret error,,,,,,err=", err)
				t.Errorf("Do error: %v", err)
			}
			// fmt.Printf("ret=%s\n", ret)
			wg.Done()
		}()
	} // --- end for ---
	time.Sleep(time.Second)
	c <- "bar"
	wg.Wait()
	fmt.Printf("count=%d\n", atomic.LoadInt32(&calls))
}

func TestExclusiveCallDoDiffDupSuppress(t *testing.T) {
	g := NewSharedCalls()
	broadcast := make(chan struct{})

	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		time.Sleep(10 * time.Millisecond)
		return nil, nil
	}
	tests := []string{"e", "a", "e", "a", "b", "c", "b", "a", "c", "d", "b", "c", "d", "d", "e"}

	var wg sync.WaitGroup
	for _, key := range tests {
		wg.Add(1)
		go func(k string) {
			<-broadcast
			_, err := g.Do(k, fn)
			if err != nil {
				t.Errorf("Do error: %v\n", err)
			}
			wg.Done()
		}(key)
	}

	time.Sleep(time.Second)
	close(broadcast)
	wg.Wait()
	fmt.Println("calls=", calls)
}

func TestExclusiveCallDoExDupSuppress(t *testing.T) {

	g := NewSharedCalls()
	c := make(chan string)

	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var wg sync.WaitGroup
	var sucCount, failCount int32
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, fresh, err := g.DoEx("key", fn)
			if err != nil {
				t.Errorf("doex failed, err=%+v\n", err)
			}
			if fresh {
				atomic.AddInt32(&sucCount, 1)
			} else {
				atomic.AddInt32(&failCount, 1)
			}
			if v.(string) != "bar" {
				t.Errorf("val compare failed, got=%+v want=%s\n", v, "bar")
			}
			wg.Done()
		}()
	}

	time.Sleep(2 * time.Second) // let goroutines above block
	c <- "bar"
	wg.Wait()
	fmt.Println("calls=", calls, ", sucCount=", sucCount, ", failCount=", failCount)
}
