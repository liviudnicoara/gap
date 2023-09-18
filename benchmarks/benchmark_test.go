package benchmarks

import (
	"sync"
	"testing"
	"time"

	"github.com/liviudnicoara/gap"
	"golang.org/x/sync/errgroup"
)

const (
	RunTimes           = 1e6
	RoutinesNumber     = 10000
	BenchParam         = 10
	DefaultExpiredTime = 10 * time.Second
)

func demoFunc(a, b int) (int, error) {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
	c := a + b
	return c, nil
}

func BenchmarkGoroutines(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			a := i
			b := j
			c := 0
			go func() {
				c, _ = demoFunc(a, b)
				wg.Done()
			}()
			c = c + 1
		}
		wg.Wait()
	}
}

func BenchmarkChannel(b *testing.B) {
	var wg sync.WaitGroup
	sema := make(chan struct{}, RoutinesNumber)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			a := i
			b := j
			c := 0
			go func() {
				c, _ = demoFunc(a, b)
				<-sema
				wg.Done()
			}()

			c = c + 1
		}
		wg.Wait()
	}
}

func BenchmarkErrGroup(b *testing.B) {
	var wg sync.WaitGroup
	var pool errgroup.Group
	pool.SetLimit(RoutinesNumber)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			a := i
			b := j
			c := 0
			pool.Go(func() error {
				c, _ = demoFunc(a, b)
				wg.Done()
				return nil
			})

			c = c + 1
		}
		wg.Wait()
	}
}

func BenchmarkGAP(b *testing.B) {
	config := gap.Config{
		BaseWorkers:   RoutinesNumber,
		MaxWorkers:    0,
		WorkerTimeout: DefaultExpiredTime,
	}
	p := gap.NewTaskPool(&config)
	defer p.Stop()

	b.ResetTimer()
	g := gap.NewGroupInPool(p)
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			a := i
			b := j
			g.Do(func() (interface{}, error) {
				return demoFunc(a, b)
			})
		}
		_ = g.GetResults()
	}
}
