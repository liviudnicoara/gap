package gap

import (
	"time"
)

// TaskPool is an interface for managing and executing tasks.
type TaskPool interface {
	Do(Task)
	Running() int
	Stop()
}

// taskPool is an implementation of the TaskPool interface.
type taskPool struct {
	baseWorkerCount        int
	maxWorkerCount         int
	temporaryWorkerPool    chan struct{}
	temporaryWorkerTimeout time.Duration
	tasks                  chan Task
	done                   chan struct{}
}

// NewTaskPool creates a new TaskPool using environment configuration.
func NewTaskPool(config *Config) TaskPool {
	base := config.BaseWorkers
	max := config.MaxWorkers
	timeout := config.WorkerTimeout

	if base <= 0 || max < 0 {
		panic("BaseWorkers and MaxWorkers must be positive values")
	}

	temporaryWorkerCount := 0
	if max > base {
		temporaryWorkerCount = max - base
	}

	appPool := &taskPool{
		baseWorkerCount:        base,
		maxWorkerCount:         max,
		temporaryWorkerPool:    make(chan struct{}, temporaryWorkerCount),
		temporaryWorkerTimeout: timeout,
		tasks:                  make(chan Task),
		done:                   make(chan struct{}),
	}

	started := make(chan struct{}, base)
	for i := 0; i < base; i++ {
		started <- struct{}{}
	}

	// Start base worker goroutines.
	for i := 0; i < base; i++ {
		w := NewTaskWorker(appPool.done, appPool.tasks)
		w.Start(started)
	}

	for i := 0; i < base; i++ {
		started <- struct{}{}
	}

	return appPool
}

// Do adds a task to the TaskPool for execution.
func (tp *taskPool) Do(task Task) {
	select {
	case tp.tasks <- task:
		return
	case tp.temporaryWorkerPool <- struct{}{}:
		// Start a temporary worker to handle the task.
		w := NewTaskWorker(tp.done, tp.tasks)
		w.StartTemporary(tp.temporaryWorkerPool, tp.temporaryWorkerTimeout)
		tp.tasks <- task
	}
}

// Running retrunrs the number of current running go routiens
func (tp *taskPool) Running() int {
	return tp.baseWorkerCount + len(tp.temporaryWorkerPool)
}

// Stop stops all workers in the TaskPool and releases associated resources.
func (tp *taskPool) Stop() {
	close(tp.done)
	close(tp.temporaryWorkerPool)
}
