package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
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
