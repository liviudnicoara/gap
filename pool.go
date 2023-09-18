package gap

import (
	"sync"
	"time"

	"github.com/liviudnicoara/gap/util"
)

// once is used for ensuring that the task pool is initialized only once.
var once sync.Once

// appPool is a singleton instance of the task pool.
var appPool *taskPool

// init initializes environment configurations using the util package.
func init() {
	util.InitEnvConfigs()
}

// TaskPool is an interface for managing and executing tasks.
type TaskPool interface {
	Do(Task)
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
func NewTaskPool() TaskPool {
	once.Do(func() {
		// Retrieve task pool configuration from environment variables using the util package.
		base := util.EnvConfigs.BaseWorkers
		max := util.EnvConfigs.MaxWorkers
		timeout := util.EnvConfigs.WorkerTimeoutInSeconds

		if base <= 0 || max <= 0 {
			panic("BaseWorkers and MaxWorkers must be positive values")
		}

		temporaryWorkerCount := 0
		if max > base {
			temporaryWorkerCount = max - base
		}

		appPool = &taskPool{
			baseWorkerCount:        base,
			maxWorkerCount:         max,
			temporaryWorkerPool:    make(chan struct{}, temporaryWorkerCount),
			temporaryWorkerTimeout: timeout,
			tasks:                  make(chan Task),
			done:                   make(chan struct{}),
		}

		// Start base worker goroutines.
		for i := 0; i < base; i++ {
			w := NewTaskWorker(appPool.done, appPool.tasks)
			w.Start()
		}
	})

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

// Stop stops all workers in the TaskPool and releases associated resources.
func (tp *taskPool) Stop() {
	close(tp.done)
	close(tp.temporaryWorkerPool)
}
