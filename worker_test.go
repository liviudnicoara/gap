package gap

import (
	"testing"
	"time"
)

func TestNewTaskWorker(t *testing.T) {
	// Create channels for testing
	done := make(chan struct{})
	tasks := make(chan Task)

	// Create a task worker
	worker := NewTaskWorker(done, tasks)

	// Check if the worker is not nil
	if worker == nil {
		t.Errorf("Expected NewTaskWorker to return a non-nil worker, but got nil")
	}
}

func TestTaskWorker_Start(t *testing.T) {

	// Create channels for testing
	done := make(chan struct{})
	tasks := make(chan Task)
	resultChannel := make(chan TaskResult)

	// Create a task worker
	worker := NewTaskWorker(done, tasks)

	// Start the task worker
	go worker.Start()

	// Create a sample task
	task := Task{
		Result: resultChannel,
		fn: func() (interface{}, error) {
			return "Hello, World!", nil
		},
	}

	// Send the task to the worker
	tasks <- task

	// Wait for the result
	result := <-resultChannel

	// Check if the task was successful
	if !result.Success {
		t.Errorf("Expected task to be successful, but it failed")
	}

	// Check the result value
	if result.Result != "Hello, World!" {
		t.Errorf("Expected result to be 'Hello, World!', but got '%v'", result.Result)
	}

	// Close the done channel to stop the worker
	close(done)
}

func TestTaskWorkerStartTemporary(t *testing.T) {
	// Create channels for testing
	done := make(chan struct{})
	tasks := make(chan Task)
	temporaryWorkerPool := make(chan struct{})
	resultChannel := make(chan TaskResult)

	// Create a task worker
	worker := NewTaskWorker(done, tasks)

	// Start the task worker temporarily
	go worker.StartTemporary(temporaryWorkerPool, 1*time.Second)

	// Create a sample task
	task := Task{
		Result: resultChannel,
		fn: func() (interface{}, error) {
			return "Temporary Worker Task", nil
		},
	}

	// Send the task to the worker
	tasks <- task

	// Wait for the result
	result := <-resultChannel

	// Check if the task was successful
	if !result.Success {
		t.Errorf("Expected task to be successful, but it failed")
	}

	// Check the result value
	if result.Result != "Temporary Worker Task" {
		t.Errorf("Expected result to be 'Temporary Worker Task', but got '%v'", result.Result)
	}

	// Close the done channel to stop the worker
	close(done)
}
