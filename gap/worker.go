package gap

import (
	"time"
)

type TaskFunc func() (interface{}, error)

type TaskResult struct {
	Success bool
	Result  interface{}
	Error   error
}

type Task struct {
	Result chan<- TaskResult
	fn     TaskFunc
}

type TaskWorker interface {
	Start()
	StartTemporary(<-chan struct{}, time.Duration)
}

type taskWorker struct {
	tasks <-chan Task
	done  <-chan struct{}
}

func NewTaskWorker(done <-chan struct{}, tasks <-chan Task) TaskWorker {
	return &taskWorker{
		tasks: tasks,
		done:  done,
	}
}

func (tw *taskWorker) Start() {
	go func() {
		for {
			select {
			case task := <-tw.tasks:
				// Do work here
				value, err := task.fn()

				// Send the result back to the caller
				task.Result <- TaskResult{
					Success: err == nil,
					Result:  value,
					Error:   err,
				}
			case <-tw.done:
				return
			}
		}
	}()
}

func (tw *taskWorker) StartTemporary(temporaryWorkerPool <-chan struct{}, timeout time.Duration) {
	go func() {
		t := time.NewTimer(timeout)

		defer func() {
			t.Stop()
			<-temporaryWorkerPool
		}()

		for {
			select {
			case task := <-tw.tasks:
				// Do work here
				value, err := task.fn()

				// Send the result back to the caller
				task.Result <- TaskResult{
					Success: err == nil,
					Result:  value,
					Error:   err,
				}
			case <-tw.done:
				return
			case <-t.C:
				return
			}
		}
	}()
}
