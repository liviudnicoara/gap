package gap

import (
	"sync"
	"time"
)

var once sync.Once
var instance *taskPool

type TaskPool interface {
	Do(Task)
	Stop()
}

type taskPool struct {
	baseWorkerCount     int
	maxWorkerCount      int
	temporaryWorkerPool chan struct{}
	tasks               chan Task
	done                chan struct{}
}

func NewTaskPool(base int, max int) TaskPool {
	once.Do(func() {
		temporaryWorkerCount := 0
		if max > base {
			temporaryWorkerCount = max - base
		}

		instance = &taskPool{
			baseWorkerCount:     base,
			maxWorkerCount:      max,
			temporaryWorkerPool: make(chan struct{}, temporaryWorkerCount),
			tasks:               make(chan Task),
			done:                make(chan struct{}),
		}

		for i := 0; i < base; i++ {
			w := NewTaskWorker(instance.done, instance.tasks)
			w.Start()
		}
	})

	return instance
}

func (tp *taskPool) Do(task Task) {
	select {
	case tp.tasks <- task:
		return
	case tp.temporaryWorkerPool <- struct{}{}:
		w := NewTaskWorker(tp.done, tp.tasks)
		w.StartTemporary(tp.temporaryWorkerPool, 10*time.Second)
		tp.tasks <- task
	}

}

func (tp *taskPool) Stop() {
	close(tp.done)
	close(tp.temporaryWorkerPool)
}
