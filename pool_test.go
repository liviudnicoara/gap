package gap

import (
	"testing"
	"time"
)

func TestNewTaskPool(t *testing.T) {
	// Create a mock configuration
	mockConfig := &Config{
		BaseWorkers:   5,
		MaxWorkers:    10,
		WorkerTimeout: 5 * time.Second,
	}

	// Create a new TaskPool
	tp := NewTaskPool(mockConfig)

	// Ensure that the taskPool is not nil
	if tp == nil {
		t.Error("Expected taskPool to be non-nil, but it is nil")
	}

	// Ensure that the taskPool is initialized with the provided configuration
	// by checking the base worker count
	if tp, ok := tp.(*taskPool); ok {
		if tp.baseWorkerCount != mockConfig.BaseWorkers {
			t.Errorf("Expected baseWorkerCount to be %d, but got %d", mockConfig.BaseWorkers, tp.baseWorkerCount)
		}
	} else {
		t.Error("Expected taskPool to be of type *taskPool, but it is not")
	}
}

func TestTaskPoolDo(t *testing.T) {
	// Create a mock configuration
	mockConfig := &Config{
		BaseWorkers:   5,
		MaxWorkers:    10,
		WorkerTimeout: 5 * time.Second,
	}

	// Create a new TaskPool
	tp := NewTaskPool(mockConfig)

	// Create a mock task
	resultChan := make(chan TaskResult)

	mockTask := Task{
		Result: resultChan,
		fn: func() (interface{}, error) {
			return "Test Result", nil
		},
	}

	// Add the mock task to the task pool
	tp.Do(mockTask)

	resultTask := <-resultChan
	result := resultTask.Result.(string)

	if result != "Test Result" {
		t.Errorf("Expected 'Test Result' task in the task pool, but got %s", result)
	}
}

func TestTaskPoolRunning(t *testing.T) {
	// Create a mock configuration
	mockConfig := &Config{
		BaseWorkers:   5,
		MaxWorkers:    10,
		WorkerTimeout: 5 * time.Second,
	}

	// Create a new TaskPool
	taskPool := NewTaskPool(mockConfig)

	// Ensure that the Running method returns the correct number of running workers
	if runningCount := taskPool.Running(); runningCount != mockConfig.BaseWorkers {
		t.Errorf("Expected %d running workers, but got %d", mockConfig.BaseWorkers, runningCount)
	}
}

func TestTaskPoolStop(t *testing.T) {
	// Create a mock configuration
	mockConfig := &Config{
		BaseWorkers:   5,
		MaxWorkers:    10,
		WorkerTimeout: 5 * time.Second,
	}

	// Create a new TaskPool
	tp := NewTaskPool(mockConfig)

	// Call the Stop method on the task pool
	tp.Stop()

	// Ensure that the done and temporaryWorkerPool channels are closed
	if tp, ok := tp.(*taskPool); ok {
		select {
		case _, open := <-tp.done:
			if open {
				t.Error("Expected done channel to be closed, but it is open")
			}
		default:
			t.Error("Expected done channel to be closed, but it is not")
		}

		select {
		case _, open := <-tp.temporaryWorkerPool:
			if open {
				t.Error("Expected temporaryWorkerPool channel to be closed, but it is open")
			}
		default:
			t.Error("Expected temporaryWorkerPool channel to be closed, but it is not")
		}
	} else {
		t.Error("Expected taskPool to be of type *taskPool, but it is not")
	}
}
