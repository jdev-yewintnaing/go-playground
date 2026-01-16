# Go Playground

This repository is a collection of small Go projects and experiments created for personal learning and practice. It covers various Go concepts, from basic syntax to concurrency and design patterns.

## Project Structure

The project follows a standard Go directory structure:

- `main.go`: The entry point of the application.
- `pkg/`: Contains reusable packages.
  - `logger/`: A custom logging library demonstrating the **Functional Options** pattern and **Interfaces**.
  - `scheduler/`: A simple task scheduler using **Goroutines** and **Channels** for concurrent execution.
  - `watcher/`: An HTTP status checker that explores **Context (with timeout)** and **Concurrency** patterns.

## Learnings & Key Concepts

### 1. Functional Options Pattern (`pkg/logger`)
- Implementing flexible configuration for structs using the functional options pattern.
- Defining interfaces to decouple the logger from its output destination (e.g., `ConsoleWriter`, `MockWriter`).
- Custom error handling with a struct that implements the `error` interface.

### 2. Concurrency with Goroutines and Channels (`pkg/scheduler`)
- Decoupling tasks using a `Task` interface.
- Running multiple tasks concurrently using `go` routines.
- Synchronizing and collecting results using channels.

### 3. Context and HTTP (`pkg/watcher`)
- Making HTTP requests that respect `context.Context` for timeouts and cancellation.
- Implementing a "Watch" function that fans out multiple requests and collects results.
- Writing table-driven tests to verify behavior under different conditions (success, timeout).

## Getting Started

### Prerequisites
- Go 1.25 or later

### Running the Application
To run the main entry point:
```bash
go run main.go
```

### Running Tests
To run all tests in the project:
```bash
go test ./...
```

To run tests with verbose output:
```bash
go test -v ./...
```
