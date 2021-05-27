package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

func dummy_somefun() {
	fmt.Println("some fun...")
}

func SearchTarget(ctx context.Context, data []int, target int, resultChan chan bool) {

	for _, v := range data {
		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stdout, "Task cancel!!!\n")
			return
		default:
		}

		// 模拟耗时查找
		fmt.Fprintf(os.Stdout, "input v: %d\n", v)
		time.Sleep(time.Millisecond * 1500)

		if target == v {
			fmt.Fprintf(os.Stdout, " found v=%d\n", target)
			resultChan <- true
			return
		}
	}

}

func tstsomeCtx() {
	timer := time.NewTimer(time.Second * 5)
	data := []int{1, 2, 3, 10, 999, 8, 345, 7, 98, 33, 66, 77, 88, 68, 96}
	dataLen := len(data)
	size := 3
	target := 345

	ctx, cancel := context.WithCancel(context.Background())
	resultChan := make(chan bool)

	for i := 0; i < dataLen; i += size {
		endPos := i + size
		if endPos >= dataLen {
			endPos = dataLen - 1
		}

		go SearchTarget(ctx, data[i:endPos], target, resultChan)
	}

	select {
	case <-timer.C:
		fmt.Fprintln(os.Stderr, "Timeout! Not Found")
		cancel()
	case <-resultChan:
		fmt.Fprintf(os.Stderr, "Found it!!!\n")
		cancel()
	}
	time.Sleep(time.Second * 2)
}

func ThreadFun() {
	time.Sleep(time.Millisecond * 800)
	fmt.Println("worker done ...")
}

func BarrFun() {
	defer fmt.Println("all done ...")

	{
		StartRunGun := make(chan struct{})

		go func() {
			time.Sleep(time.Second * 2)

			fmt.Println("aaaaaaaaaaaaaaaa")
			close(StartRunGun)
			fmt.Println("bbbbbbbbbbbbbbbb")
		}()
		<-StartRunGun
		fmt.Println("cccccccccccccccccc")
		return
	}
	{
		fmt.Println("Hello World!")
		barrier := NewBarrier(3)
		var wg sync.WaitGroup

		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println("A")
				barrier.BarrierWait()
				fmt.Println("B")
				barrier.BarrierWait()
				fmt.Println("C")
			}()
		}

		wg.Wait()
		return
	}

	var count int = 5
	barrobj := NewBarrier(count)

	for i := 0; i < count; i++ {
		go ThreadFun()
	}

	barrobj.BarrierWait()

}
