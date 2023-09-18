package gap

import (
	"fmt"
	"time"
)

// TaskFunc represents a function that can be executed as a task.
type TaskFunc func() (interface{}, error)

// TaskResult holds the result of a task execution.
type TaskResult struct {
	Success bool        // Indicates whether the task was successful.
	Result  interface{} // The result of the task (if successful).
	Error   error       // An error (if the task encountered an error).
}

// Task represents a unit of work that can be executed by a TaskWorker.
type Task struct {
	Result chan<- TaskResult // A channel to send the task result to the caller.
	fn     TaskFunc          // The function to be executed as the task.
}

// TaskWorker is an interface for executing tasks.
type TaskWorker interface {
	Start()
	StartTemporary(<-chan struct{}, time.Duration)
}

// taskWorker is an implementation of the TaskWorker interface.
type taskWorker struct {
	tasks <-chan Task     // A channel to receive tasks.
	done  <-chan struct{} // A channel to signal when the worker should stop.
}

// NewTaskWorker creates a new TaskWorker.
func NewTaskWorker(done <-chan struct{}, tasks <-chan Task) TaskWorker {
	return &taskWorker{
		tasks: tasks,
		done:  done,
	}
}

// Start starts the task worker, continuously executing tasks.
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

// StartTemporary starts a temporary task worker with a timeout.
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
				fmt.Println("timeout")
				return
			}
		}
	}()
}
