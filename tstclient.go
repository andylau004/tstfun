package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"utils"

	"go.uber.org/zap"
)

func getResponse(resp *http.Response, err error, id int) {

	time.Sleep(5 * time.Second)

	buf, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%d: %s -- %v\n", id, string(buf), err)

	if err := resp.Body.Close(); err != nil {
		fmt.Println(err)
	}

}

func doGet(client *http.Client, url string, id int) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%d: %s -- %v\n", id, string(buf), err)

	if err := resp.Body.Close(); err != nil {
		fmt.Println(err)
	}
	return

	time.Sleep(51 * time.Second)
	go getResponse(resp, err, id)
}

const URL = "http://192.168.6.127:3008/file-agent/fileagent/test"

func tstHttpClient1() {

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: PrintLocalDial,
		},
	}
	for {
		go doGet(httpClient, URL, 1)

		// go doGet(httpClient, URL, 2)
		break
		time.Sleep(2 * time.Second)
	}
	// time.Sleep(2 * time.Second)
	select {}
}

func PrintLocalDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	conn, err := dial.Dial(network, addr)
	if err != nil {
		return conn, err
	}

	fmt.Println("connect done, use", conn.LocalAddr().String())
	return conn, err
}

var gIdx int32

func onTest(w http.ResponseWriter, req *http.Request) { // 中孚会到我们后置服务处，请求豆豆token，此接口处理

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		gLogger.Error("found bd login failed!", zap.Error(err))
		return
	}
	defer req.Body.Close()

	_ = bodyBytes

	strRet := fmt.Sprintf("%d_resp to client", atomic.AddInt32(&gIdx, 1))
	w.Write(utils.Str2bytes(strRet))
}

func tstHttpSrv() {
	fmt.Println("----------------beg http")
	defer fmt.Println("----------------end http")

	hostPort := "0.0.0.0:54321"
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", onTest)

	server := &http.Server{
		Addr:    hostPort,
		Handler: mux,
		// WriteTimeout: time.Second * 3, //设置3秒的写超时
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("new bd server listen failed! err=", err)
		// gLogger.Error("new bd server listen failed! err=", zap.Error(err), zap.String("listAddr=", hostPort))
		panic(err)
	}
}

func tstHttpLongConnClient() {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: false,
	}
	client := &http.Client{Transport: tr}

	var wg sync.WaitGroup //用来在后面同步所有协程结束
	for i := 0; i < 10; i++ {
		wg.Add(1) //增加计数器

		go func() {

			defer wg.Done() //减少计数器

			resp, err := client.Get("http://192.168.6.64:54321/hello")
			if err != nil {
				fmt.Println("get req failed, err=", err)
				return
			}

			slurp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("read resp failed, err=", err)
				return
			}

			fmt.Printf("got resp=%+v\n", string(slurp))
			resp.Body.Close()

			select {}
			// time.Sleep(time.Millisecond * 500)
			// time.Sleep(time.Millisecond * 500)
			// time.Sleep(time.Millisecond * 500)
		}()

	}
	wg.Wait() //等待所有协程结束

}
