package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func doTimeStuff(ctx context.Context) {

	for {
		time.Sleep(time.Second)

		if deadline, ok := ctx.Deadline(); ok {
			fmt.Println("deadline have been set")

			if time.Now().After(deadline) {
				fmt.Printf(" after deadline, err=%+v\n", ctx.Err().Error())
				return
			}
		}

		select {
		case <-ctx.Done():
			fmt.Println("done ...")
		default:
			fmt.Println("work ...")
		}

	} // end for

}

func tstCtxFun1() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go doTimeStuff(ctx)

	time.Sleep(10 * time.Second)

	fmt.Println("beg exe cancel")
	cancel()
	fmt.Println("end exe cancel")
}

func tstCtxFun2() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	for {

		select {
		case <-time.After(1 * time.Second):
			fmt.Println("overslept")
		case <-ctx.Done():
			fmt.Println(ctx.Err()) // prints "context deadline exceeded"
		}

	}

}

func tstCtxFun3() {
	uri := "https://httpbin.orgttt/delay/3"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Printf("http.NewRequest() failed with '%s'\n", err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*4)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("http.DefaultClient.Do() failed with:\n'%s'\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("body=[%s]\n", string(body))
	fmt.Println("do req done...")
}
func tstCtxFun56() {
	uri := "http://192.168.6.129:3008/file-agent/fileagent/test"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Printf("http.NewRequest() failed with '%s'\n", err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*4)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("http.DefaultClient.Do() failed with:\n'%s'\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("body=[%s]\n", string(body))
	fmt.Println("do req done...")
}

func func1(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	respC := make(chan int)

	go func() { // work thread code
		time.Sleep(1 * time.Second)
		// respC <- 13
	}()

	select {
	case <-ctx.Done():
		fmt.Println("ctx done...")
	case val := <-respC:
		fmt.Println("val=", val)
	}

	return nil
}

func tstCtxFun4() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	go func1(ctx, wg)

	time.Sleep(3 * time.Second)

	cancel()
	wg.Wait()
}

func func15(ctx context.Context) {
	ctx = context.WithValue(ctx, "k1", "v1")
	func25(ctx)
}
func func25(ctx context.Context) {
	fmt.Println(ctx.Value("k1").(string))
}
func tstCtxFun5() {
	ctx := context.Background()
	func15(ctx)
}

func tstCtxFun() {
	tstCtxFun4()
	return

	tstCtxFun56()
	return

	tstCtxFun3()
}
