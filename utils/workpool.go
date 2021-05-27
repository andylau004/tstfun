package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

type TaskFuncType func(i int)

type Task struct {
	TaskFunc TaskFuncType
	Idx      int
	// TaskFunc func(workID int, threadData interface{})
}

type WorkerPool struct {
	worker_num  int // 启动的协程个数
	task_num    int
	jobsChannel chan *Task
	chDone      chan int
	// chStart     chan struct{}
	barrier  *Barrier
	stop_num int32
}

//创建一个协程池
// func NewWorkerPool(workerNum, tasknum int, chStart chan struct{}) *WorkerPool {
func NewWorkerPool(workerNum, tasknum int) *WorkerPool {
	maxWait := workerNum
	if workerNum < tasknum {
		maxWait = tasknum
	}
	maxWait += 1

	p := WorkerPool{
		worker_num:  workerNum,
		task_num:    tasknum,
		jobsChannel: make(chan *Task, maxWait),
		chDone:      make(chan int),
		// chStart:     chStart,
		barrier:  NewBarrier(workerNum),
		stop_num: 0,
	}
	return &p
}

func (p *WorkerPool) AddTask(t *Task) {
	p.jobsChannel <- t
}

func (p *WorkerPool) StopOneWorker() {
	newstopnum := atomic.AddInt32(&p.stop_num, 1)
	if newstopnum == int32(p.task_num) {
		close(p.jobsChannel)
	}
}

var onceSetTimeBegin sync.Once

//协程池创建一个worker并且开始工作
func (p *WorkerPool) worker(wg *sync.WaitGroup, workID int, timeBegin *time.Time) {
	defer wg.Done()

	p.barrier.BarrierWait()

	// <-p.chStart

	onceSetTimeBegin.Do(func() {
		*timeBegin = time.Now()
	})

	var isExit bool = false
	// threadData := onInit()
	// fmt.Println("workerID=", workID, " already ...")
	for {
		select {
		case task, isOK := <-p.jobsChannel:
			if !isOK {
				isExit = true
				break
			}

			task.TaskFunc(task.Idx)
			// task.TaskFunc(workID, threadData)
		case <-p.chDone:
			isExit = true
			break
		}
		if isExit {
			break
		}
	}

	// onDestroy(threadData)
}

func (p *WorkerPool) CommitTaskOver() {
	close(p.chDone)
}

//让协程池Pool开始工作
func (p *WorkerPool) Run(wg *sync.WaitGroup, timeBegin *time.Time) {
	for i := 0; i < p.worker_num; i++ {
		go p.worker(wg, i, timeBegin)
	}
}

// func (p *WorkerPool) Run(wg *sync.WaitGroup, timeBegin *time.Time, onInit func() interface{}, onDestroy func(interface{})) {
// 	for i := 0; i < p.worker_num; i++ {
// 		go p.worker(wg, i, timeBegin, onInit, onDestroy)
// 	}
// }
