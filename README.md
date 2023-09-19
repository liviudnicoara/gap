![](https://raw.githubusercontent.com/liviudnicoara/liviudnicoara/main/assests/pool.png)

# GAP - A Global Go App Pool

## 📖 Introduction

Library `gap` implements a goroutine pool with initial fixed capacity, allowing developers to limit the number of goroutines in your concurrent programs.


## 🛠 Use cases
GAP was designed to provided a limitation on spanned go routines in a business flow. For example, GAP can be used to limit all http calls that you application does (see Usage section below)

## 🚀 Features:

- **Task Pool Management**: Go App Pool allows you to create and manage a task pool with configurable base worker counts, max worker counts and timeout settings.

- **Task Grouping**: You can use the provided Group to manage and execute a group of tasks concurrently and collect their results.

- **Temporary workers**: You can configure GAP to create the desired number of workers during a high load and to stop them after an idle timeout.

- **Environment Variable Configuration**: Go App Pool supports configuration through environment variables, making it easy to adjust worker counts and timeouts without modifying the code.


## 💡 How `gap` works

1. When your application starts, gap creates the base number of go routines defined in the configuration
2. The active workers will listent for any new task
3. The developer creates a group of tasks and asigns a function to each task using the `Do` method (simmilar to errorgroup library)
4. Afterwards, it waits for the results using `GetResults` method
5. If the number of tasks exceeds the base number of goroutines, `gap` creates temprary workers up to the max value defined in the config.
6. After an idle tiemoput time, the temprary workers are stopped.


## 🧰 How to install

``` powershell
go get github.com/liviudnicoara/gap
```

```go
import "github.com/liviudnicoara/gap"
```

### Usage
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

`GAP_BASE_WORKERS`: The number of base worker goroutines that will be active while your application is running. Defaults to `100`.

`GAP_MAX_WORKERS`: The maximum number of worker goroutines. Temporary workers will be created during high load and removed after an idle time. Defaults to `1000`.

`GAP_WORKER_TIMEOUT`: The timeout for temporary worker goroutines  Defaults to `1s`.

### Examples
For code examples and advanced usage, please refer to the examples directory.

### Contributing
Contributions are welcome! If you'd like to contribute to Go App Pool, please follow the contribution guidelines.

### License
This project is licensed under the MIT License - see the LICENSE file for details.

