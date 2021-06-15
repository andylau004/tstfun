package taskpool

import "sync"

type TaskFn func(params ...interface{})

type taskWrapper struct {
	taskFn TaskFn
	param  []interface{}
}

type Status struct {
	TotalWorkerNum int // 总协程数量
	IdleWorkerNum  int // 空闲协程数量
	BlockTaskNum   int // 等待执行的任务数。注意，只在协程数量有最大限制的情况下，该值才可能不为0，具体见Option.MaxWorkerNum
}

type Pool interface {
	// 向池内放入任务
	// 非阻塞函数，不会等待task执行
	Go(task TaskFn, param ...interface{})

	// 获取当前的状态，注意，只是一个瞬时值
	GetCurrentStatus() Status

	// 关闭池内所有的空闲协程
	KillIdleWorkers()
}

type pool struct {
	maxWorkerNum   int
	m              sync.Mutex
	totalWorkerNum int
	idleWorkerList []*worker
	blockTaskList  []taskWrapper
}

func (p *pool) Go(task TaskFn, param ...interface{}) {
	tw := &taskWrapper{
		taskFn: task,
		param:  param,
	}

}
func (p *pool) KillIdleWorkers() {
}

func (p *pool) GetCurrentStatus() Status {
	p.m.Lock()
	defer p.m.Unlock()
	return Status{
		TotalWorkerNum: p.totalWorkerNum,
		IdleWorkerNum:  len(p.idleWorkerList),
		BlockTaskNum:   len(p.blockTaskList),
	}
}
func (p *pool) newWorker() *worker {
}
func (p *pool) newWorkerWithTask(task taskWrapper) {
}
func (p *pool) onIdle(w *worker) {
	p.m.Lock()
	defer p.m.Unlock()

	if len(p.blockTaskList) == 0 {

		// 没有等待执行的任务
		p.idleWorkerList = append(p.idleWorkerList, w)

	} else {
		
	}

}

func NewPool() (Pool, error) {
	p := pool{
		maxWorkerNum: 10,
	}
	return &p, nil
}

// func newPool() (*pool, error) {
// 	p := pool{
// 		maxWorkerNum: 10,
// 	}
// 	return &p, nil
// }
