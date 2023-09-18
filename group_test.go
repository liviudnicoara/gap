package gap

import (
	"testing"
)

func TestTaskGroup(t *testing.T) {
	// Create a new TaskGroup
	taskGroup := NewGroup()

	// Define a test function that simulates a task
	testTask := func() (interface{}, error) {
		return "Test Result", nil
	}

	// Add a task to the TaskGroup
	taskGroup.Do(testTask)

	// Get the results from the TaskGroup
	results := taskGroup.GetResults()

	// Check if the TaskGroup contains the expected number of results (1 in this case)
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Check the result value
	if results[0].Result != "Test Result" {
		t.Errorf("Expected result 'Test Result', got '%v'", results[0].Result)
	}

}

func TestTaskGroupEmpty(t *testing.T) {
	// Create a new TaskGroup
	taskGroup := NewGroup()

	// Get the results from an empty TaskGroup
	results := taskGroup.GetResults()

	// Check if the TaskGroup contains no results
	if len(results) != 0 {
		t.Errorf("Expected 0 results for an empty TaskGroup, got %d", len(results))
	}
}

func TestTaskGroupMultipleTasks(t *testing.T) {
	// Create a new TaskGroup
	taskGroup := NewGroup()

	// Define a test function that simulates a task
	testTask := func() (interface{}, error) {
		return "Task Result", nil
	}

	// Add multiple tasks to the TaskGroup
	for i := 0; i < 5; i++ {
		taskGroup.Do(testTask)
	}

	// Get the results from the TaskGroup
	results := taskGroup.GetResults()

	// Check if the TaskGroup contains the expected number of results (5 in this case)
	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}
}
