package gap

import (
	"sync"
)

// TaskGroup represents a group of tasks executed concurrently using a task pool.
type TaskGroup struct {
	taskPool       TaskPool        // taskPool is the pool of workers for task execution.
	wg             sync.WaitGroup  // wg is used to wait for all tasks to complete.
	resultsChannel chan TaskResult // resultsChannel is used to collect task results.
	resultsStore   []TaskResult    // resultsStore stores the collected task results.
}

// NewGroup creates a new TaskGroup for managing and executing a group of tasks concurrently.
func NewGroup() *TaskGroup {
	return NewGroupInPool(defaultPool)
}

// NewGroup creates a new TaskGroup for managing and executing a group of tasks concurrentlyin a specified pool.
func NewGroupInPool(pool TaskPool) *TaskGroup {
	taskGroup := &TaskGroup{
		taskPool:       pool,                  // Create a task pool with default configuration.
		resultsChannel: make(chan TaskResult), // Create a channel to collect task results.
		resultsStore:   []TaskResult{},        // Initialize an empty slice to store task results.
	}

	// Start a goroutine that will gather results
	taskGroup.dispatchResults()

	return taskGroup
}

func (tg *TaskGroup) dispatchResults() {
	// Start a goroutine that will gather results
	sem := make(chan struct{}, 1)
	defer close(sem)

	sem <- struct{}{}

	go func(started <-chan struct{}) {
		<-started

		for r := range tg.resultsChannel {
			tg.resultsStore = append(tg.resultsStore, r)
			tg.wg.Done()
		}
	}(sem)

	sem <- struct{}{}
}

// Do adds a task to the TaskGroup for execution.
func (tg *TaskGroup) Do(fn func() (interface{}, error)) {
	tg.wg.Add(1)
	tg.taskPool.Do(Task{
		Result: tg.resultsChannel,
		fn:     fn,
	})
}

// GetResults waits for all tasks in the TaskGroup to complete and returns the collected task results.
func (g *TaskGroup) GetResults() []TaskResult {
	g.wg.Wait()
	close(g.resultsChannel)
	return g.resultsStore
}
