package utils

import (
	"sync"
	"testing"
	"time"
)

// "github.com/pquerna/ffjson/ffjson"

func Test_threadpool_1(t *testing.T) {

	var tmBeg time.Time
	var wg sync.WaitGroup

	var threadNUm int = 5
	workerPool := NewWorkerPool(threadNUm, threadNUm)

	wg.Add(threadNUm)

	workerPool.Run(&wg, &tmBeg)

	wg.Wait()
}
