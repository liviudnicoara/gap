![](https://raw.githubusercontent.com/liviudnicoara/liviudnicoara/main/assests/pool.png)

# GAP - A Global Go App Pool

GAP is a GoLang library that provides a task pool for managing and executing tasks concurrently. It allows you to efficiently execute multiple tasks using a pool of worker goroutines, making it suitable for tasks that can be parallelized.

## Usecases
GAP was designed to provided a limitation on spanned go routines in a business flow. For example, GAP can be used to limit all http calls that you application does (see Usage section below)

## Features

- **Task Pool Management**: Go App Pool allows you to create and manage a task pool with configurable base worker counts, max worker counts and timeout settings.

- **Task Grouping**: You can use the provided Group to manage and execute a group of tasks concurrently and collect their results.

- **Temporary workers**: You can configure GAP to create the desired number of workers during a high load and to stop them after an idle timeout.

- **Environment Variable Configuration**: Go App Pool supports configuration through environment variables, making it easy to adjust worker counts and timeouts without modifying the code.

## Getting Started

### Prerequisites

- Go (Golang) must be installed on your system. You can download it from [here](https://golang.org/dl/).

### Installation

1. Clone the repository:

   ```shell
   git clone https://github.com/liviudnicoara/gap
   ```

2. Navigate to the project directory:

    ```shell
    cd go-app-pool
    ```

3. Build the application:

    ```shell
    go build
    ```

4. Run the example:
    ```shell
    go run .
    ```

### Usage
Import the gap package in your Go application:

```go
import "github.com/liviudnicoara/gap"
```
Create a TaskPool using the provided configuration:

```go
taskPool := gap.NewTaskPool()
```
Use the TaskGroup to manage and execute tasks concurrently:

```go

taskGroup := gap.NewGroup()

// Add tasks to the task group
taskGroup.Do(func() (interface{}, error) {
    // Your task logic here
})

// Wait for all tasks to complete and collect results
results := taskGroup.GetResults()
```
Create task groups for different uses case: 

```go

defer gap.Stop()

	alphaCodes := []string{	"USA", "CAN", "GBR", "FRA" }

	countryGroup := gap.NewGroup()

	for _, c := range alphaCodes {
		code := c
		countryGroup.Do(func() (interface{}, error) {
			return GetCommonNameByAlphaCode(code)
		})
	}

	todoGroup := gap.NewGroup()

	fmt.Println("Active go routines: ", gap.Running())

	for i := 1; i < 5; i++ {
		id := i
		countryGroup.Do(func() (interface{}, error) {
			return GetTodoByID(id)
		})
	}

	fmt.Println("Getting country results")
	countryResults := todoGroup.GetResults()
	for _, r := range countryResults {
		fmt.Println(r.Result)
	}

	fmt.Println("Getting todo results")
	todoResults := countryGroup.GetResults()
	for _, r := range todoResults {
		fmt.Println(r.Result)
	}
```

Customize the task pool configuration by adjusting environment variables (see Configuration section below).

### Configuration
Go App Pool can be configured using environment variables. Here are the available configuration options:

GAP_BASE_WORKERS: The number of base worker goroutines that will be active while your application is running

GAP_MAX_WORKERS: The maximum number of worker goroutines. Temporary workers will be created during high load and removed after an idle time

GAP_WORKER_TIMEOUT: The timeout for temporary worker goroutines (default: "10s").

### Examples
For code examples and advanced usage, please refer to the examples directory.

### Contributing
Contributions are welcome! If you'd like to contribute to Go App Pool, please follow the contribution guidelines.

### License
This project is licensed under the MIT License - see the LICENSE file for details.
