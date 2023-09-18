package gap

var (
	configs = loadEnvVariables()

	defaultPool = NewTaskPool(configs)
)

// Running retrunrs the number of current running go routiens
func Running() int {
	return defaultPool.Running()
}

// Stop stops all workers in the TaskPool and releases associated resources.
func Stop() {
	defaultPool.Stop()
}
