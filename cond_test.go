package main

import (
	"log"
	"sync"
	"testing"
	"time"
)

var done = false

func read(name string, c *sync.Cond) {
	c.L.Lock()
	if !done {
		c.Wait()
	}
	log.Println(name, " start reading")
	c.L.Unlock()
}
func write(name string, c *sync.Cond) {
	log.Print(name, " start writing")
	// time.Sleep(time.Second)

	c.L.Lock()
	done = true
	c.L.Unlock()

	log.Print(name, " wakes all...")
	c.Broadcast()
}
func Test_cond1(t *testing.T) {

	cond := sync.NewCond(&sync.Mutex{})

	go read("read1", cond)
	go read("read2", cond)
	go read("read3", cond)
	time.Sleep(time.Second * 1)

	write("write1", cond)

	time.Sleep(time.Second * 3)

	_ = t
}
