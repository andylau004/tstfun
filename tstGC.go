package main

import ()
	"fmt"
)

func tstGcFun() {
	const N = 5e7 // 5000w

	if len(os.Args) != 2 {
	  fmt.Printf("usage: %s [1 2 3 4]\n(number selects the test)\n", os.Args[0])
	  return
	}
  

	
	switch os.Args[1] {
	case "1":
		m := make(map[int32]*int32)
		for i:= 0; i < N; i ++ {
			var val int32 = int32(i)
			m[val] = &val
		}
	case "2":

	}

}