https://raw.githubusercontent.com/liviudnicoara/liviudnicoara/main/assests/pool.png

# GAP - A Global Go App Pool

GAP is a GoLang library that provides a task pool for managing and executing tasks concurrently. It allows you to efficiently execute multiple tasks using a pool of worker goroutines, making it suitable for tasks that can be parallelized.

## Features

- **Task Pool Management**: Go App Pool allows you to create and manage a task pool with configurable worker counts and timeout settings.

- **Task Grouping**: You can use the provided TaskGroup to manage and execute a group of tasks concurrently and collect their results.

- **Environment Variable Configuration**: Go App Pool supports configuration through environment variables, making it easy to adjust worker counts and timeouts without modifying the code.

## Getting Started

### Prerequisites

- Go (Golang) must be installed on your system. You can download it from [here](https://golang.org/dl/).

### Installation

1. Clone the repository:

   ```shell
   git clone <repository-url>
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
import "github.com/your-username/go-app-pool/gap"
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

Customize the task pool configuration by adjusting environment variables (see Configuration section below).

### Configuration
Go App Pool can be configured using environment variables. Here are the available configuration options:

BASE_WORKERS_ENV: The number of base worker goroutines (default: 10).

MAX_WORKERS_ENV: The maximum number of worker goroutines (default: 10).

WORKER_TIMEOUT_ENV: The timeout for temporary worker goroutines (default: "10s").

### Examples
For code examples and advanced usage, please refer to the examples directory.

### Contributing
Contributions are welcome! If you'd like to contribute to Go App Pool, please follow the contribution guidelines.

### License
This project is licensed under the MIT License - see the LICENSE file for details.
