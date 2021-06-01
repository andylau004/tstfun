package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

type Obj_st struct {
	str string
	val int32
}

func tstObj() {
	var obj Obj_st
	obj.str = "123"
	obj.val = 456
	fmt.Printf("obj type: %T\n", obj)

	var at Att
	fmt.Printf("at type: %T\n", &at)
}

func tstGcFun() {
	tstObj()
	// const N = 5e6 // 5000w
	const N = 500 // 5000w
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

	case "3":
		// Split the map into 100 shards
		shards := make([]map[int32]*int32, 100)
		for i := range shards {
			shards[i] = make(map[int32]*int32)
		}
		for i := 0; i < N; i++ {
			n := int32(i)
			shards[i%100][n] = &n
		}

		// for i := range shards {
		// 	var output string
		// 	for _, v := range shards[i] {
		// 		output += fmt.Sprintf(" %d ", *v)
		// 	}
		// 	fmt.Printf("map_%d output=%s\n", i, output)
		// }

		runtime.GC()
		fmt.Printf("With map shards (%T), GC took %s\n", shards, timeGC())
	case "4":
		shards := make([]map[int32]int32, 100) // Split the map into 100 shards
		for i := range shards {
			shards[i] = make(map[int32]int32)
		}
		for i := 0; i < N; i++ {
			n := int32(i)
			shards[i%100][n] = n
		}
		runtime.GC()
		fmt.Printf("With map shards (%T), GC took %s\n", shards, timeGC())
	case "5":
		// A slice, just for comparison to show that
		// merely holding onto millions of int32s is fine
		// if they're in a slice.
		type t struct {
			p, q int32
		}
		var s []t
		for i := 0; i < N; i++ {
			n := int32(i)
			s = append(s, t{n, n})
		}
		runtime.GC()
		fmt.Printf("With a plain slice (%T), GC took %s\n", s, timeGC())
		_ = s[0]
	}

}

func timeGC() time.Duration {
	start := time.Now()
	runtime.GC()
	return time.Since(start)
}
