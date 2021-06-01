package main

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type ChanMessage interface {
	Key() string
	Value() interface{}
}

type MsgString struct {
	key   string
	value string
}

func (x *MsgString) Key() string        { return x.key }
func (x *MsgString) Value() interface{} { return x.value }

func Test_channel_NewTimer(t *testing.T) {

	Convey("It should success get two messages when no timeout", t, func() {
		var response string

		channel := make(chan ChanMessage, 2)
		defer close(channel)

		go func() {
			defer func() { recover() }() // recoverint from panic caused by writing to a closed channel
			time.Sleep(50 * time.Millisecond)
			channel <- &MsgString{
				key:   "online",
				value: "online data, ",
			}
		}()
		go func() {
			defer func() { recover() }()
			time.Sleep(50 * time.Millisecond)
			channel <- &MsgString{
				key:   "offline",
				value: "offline data, ",
			}
		}()

		timeout := time.NewTimer(1000 * time.Millisecond)

	WaitLoop1:
		for i := 0; i < cap(channel); i++ {
			select {
			case msg := <-channel:
				key, value := msg.Key(), msg.Value()
				switch key {
				case "online":
					if online_str, ok := value.(string); ok {
						response += online_str
					}
				case "offline":
					if offline_str, ok := value.(string); ok {
						response += offline_str
					}
				}

			case <-timeout.C:
				t.Fatal("Channel error timeout")
				break WaitLoop1
			}
		}

		t.Log(response)
		So(response, ShouldHaveLength, 27)
	})
}
