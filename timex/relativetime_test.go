package timex

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelativeTime(t *testing.T) {
	// time.Sleep(time.Millisecond)
	now := Now()
	// assert.True(t, now > 0)
	fmt.Println(now)
	fmt.Println("")

	time.Sleep(time.Millisecond)
	// assert.True(t, Since(now) > 0)
	fmt.Println(initTime)
	fmt.Println(Since(now))
}

func TestRelativeTime_Time(t *testing.T) {
	fmt.Println(time.Now())
	fmt.Println(Time())
	diff := time.Until(Time())
	if diff > 0 {
		assert.True(t, diff < time.Second)
	} else {
		assert.True(t, -diff < time.Second)
	}
	fmt.Println(diff)

}
