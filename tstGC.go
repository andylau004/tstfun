package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func tstGcFun() {
	const N = 5e6 // 5000w
	fmt.Println("N=", N)

	if len(os.Args) != 2 {
		fmt.Printf("usage: %s [1 2 3 4]\n(number selects the test)\n", os.Args[0])
		return
	}

	switch os.Args[1] {
	case "1":
		m := make(map[int32]*int32)
		for i := 0; i < N; i++ {
			var val int32 = int32(i)
			m[val] = &val
		}
		runtime.GC()
		fmt.Printf("With %T, GC took %s\n", m, timeGC())
		// _ = m[0] // Preserve m until here, hopefully

	case "2":
		m := make(map[int32]int32)
		for i := 0; i < N; i++ {
			var val int32 = int32(i)
			m[int32(i)] = val
		}
		runtime.GC()
		fmt.Printf("With %T, GC took %s\n", m, timeGC())
	}

}

func timeGC() time.Duration {
	start := time.Now()
	runtime.GC()
	return time.Since(start)
}
