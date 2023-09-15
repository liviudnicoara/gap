package gap

import (
	"sync"
)

type TaskGroup struct {
	taskPool       TaskPool
	wg             sync.WaitGroup
	resultsChannel chan TaskResult
	resultsStore   []TaskResult
}

func NewGroup() *TaskGroup {
	taskGroup := &TaskGroup{
		taskPool:       NewTaskPool(10, 10),
		resultsChannel: make(chan TaskResult),
		resultsStore:   []TaskResult{},
	}

	// Start go routine that will gather results
	sem := make(chan struct{})
	defer close(sem)

	go func(started <-chan struct{}) {
		sem <- struct{}{}

		for r := range taskGroup.resultsChannel {
			taskGroup.resultsStore = append(taskGroup.resultsStore, r)
			taskGroup.wg.Done()
		}
	}(sem)

	<-sem

	return taskGroup
}

func (tg *TaskGroup) Do(fn func() (interface{}, error)) {
	tg.wg.Add(1)
	tg.taskPool.Do(Task{
		Result: tg.resultsChannel,
		fn:     fn,
	})
}

func (g *TaskGroup) GetResults() []TaskResult {
	g.wg.Wait()
	close(g.resultsChannel)
	return g.resultsStore
}
