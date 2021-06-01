package syncx

import "sync"

type (

	// SharedCalls lets the concurrent calls with the same key to share the call result.
	// For example, A called F, before it's done, B called F. Then B would not execute F,
	// and shared the result returned by F which called by A.
	// The calls with the same key are dependent, concurrent calls share the returned values.
	// A ------->calls F with key<------------------->returns val
	// B --------------------->calls F with key------>returns val

	SharedCalls interface {
		Do(key string) (interface{}, error)
		DoEx(key string) (interface{}, error)
	}
)
