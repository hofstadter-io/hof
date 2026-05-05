---
applyTo: **/*.go
---

## Language Conventions

- Follow Effective Go: Adhere to official Go style guidelines and idioms
- Use gofmt and goimports: Ensure all code is properly formatted before committing
- Run go vet: Check for suspicious constructs and potential issues
- Follow Go naming conventions: Use camelCase for unexported, PascalCase for exported identifiers
- Write idiomatic Go: Prefer composition over inheritance, use interfaces appropriately

## Go Project Structure

- Standard Go layout: Follow the standard Go project layout (cmd/, internal/, pkg/, etc.)
- Package organization: Keep packages focused and avoid circular dependencies
- Internal packages: Use internal/ directory for code that shouldn't be imported externally
- Go modules: Use Go modules for dependency management (go.mod/go.sum)

## Go Testing Best Practices

- Use standard testing package: Leverage Go's built-in testing framework
- Table-driven tests: Use table-driven tests for testing multiple scenarios
- Test file naming: Name test files with _test.go suffix
- Benchmark tests: Include benchmark tests for performance-critical code
- Race detection: Run tests with -race flag to detect race conditions
- Test coverage: Use go test -cover to measure test coverage

## Go-Specific Testing Commands

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

## Error Handling

- Explicit error handling: Check and handle errors explicitly, don't ignore them
- Error wrapping: Use fmt.Errorf with %w verb or errors.Wrap for error context
- Custom error types: Create custom error types when appropriate
- Sentinel errors: Use sentinel errors for expected error conditions

## Concurrency Best Practices

- Goroutine management: Always ensure goroutines can exit cleanly
- Channel usage: Use channels for communication, mutexes for shared state
- Context usage: Use context.Context for cancellation and timeouts
- Race condition prevention: Be vigilant about data races in concurrent code
- Worker pool patterns: Implement proper worker pool patterns when needed

## Performance and Memory Management

- Memory profiling: Use go tool pprof for memory and CPU profiling
- Avoid memory leaks: Be careful with goroutines, channels, and large data structures
- Efficient string operations: Use strings.Builder for string concatenation
- Slice management: Understand slice capacity and growth patterns
- Interface efficiency: Be mindful of interface allocation overhead

## Go Build and Development Tools

- Use go build tags: Leverage build tags for environment-specific code
- Vendoring: Use go mod vendor when necessary for reproducible builds
- Static analysis: Use tools like golint, golangci-lint, and staticcheck
- Documentation: Write package documentation following Go doc conventions
- Go generate: Use go generate for code generation when appropriate

## Dependency Management

- Minimal dependencies: Keep external dependencies to a minimum
- Dependency updates: Regularly update dependencies and check for security issues
- Go mod tidy: Run go mod tidy to clean up unused dependencies
- Semantic versioning: Follow semantic versioning for your own modules