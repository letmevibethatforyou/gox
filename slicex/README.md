# sliceX

A Go package providing utility functions for slice operations with full generics support.

## Features

- **Type-safe**: All functions use Go generics for compile-time type safety
- **Pure functions**: No side effects on input slices
- **Performance optimized**: Efficient algorithms with minimal allocations
- **Zero dependencies**: Uses only Go standard library

## Functions

### Unique

Returns a new slice containing only unique elements from the input slice, preserving the order of first occurrence.

```go
func Unique[T comparable](slice []T) []T
```

**Example:**
```go
numbers := []int{1, 2, 2, 3, 1, 4, 3}
unique := slicex.Unique(numbers)
// Result: [1, 2, 3, 4]

words := []string{"hello", "world", "hello", "go"}
uniqueWords := slicex.Unique(words)
// Result: ["hello", "world", "go"]
```

### FilterNonZero

Returns a new slice with all non-zero values from the input slice. Zero values are determined by Go's zero value concept (0, "", nil, etc.).

```go
func FilterNonZero[T comparable](slice []T) []T
```

**Example:**
```go
numbers := []int{1, 0, 2, 0, 3, 0, 4}
filtered := slicex.FilterNonZero(numbers)
// Result: [1, 2, 3, 4]

strings := []string{"hello", "", "world", "", "test"}
filteredStrings := slicex.FilterNonZero(strings)
// Result: ["hello", "world", "test"]
```

### Map

Applies the given function to each element of the slice and returns a new slice containing the results.

```go
func Map[T, R any](slice []T, fn func(T) R) []R
```

**Example:**
```go
numbers := []int{1, 2, 3, 4}
strings := slicex.Map(numbers, func(n int) string {
    return strconv.Itoa(n)
})
// Result: ["1", "2", "3", "4"]

words := []string{"hello", "world", "go"}
lengths := slicex.Map(words, func(s string) int {
    return len(s)
})
// Result: [5, 5, 2]
```

### Group

Groups the elements of the slice by the result of the key function. Returns a map where keys are the grouping criteria and values are slices of grouped items.

```go
func Group[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T
```

**Example:**
```go
words := []string{"hello", "world", "go", "test", "a", "b"}
grouped := slicex.Group(words, func(s string) int {
    return len(s)
})
// Result: map[int][]string{
//   1: ["a", "b"],
//   2: ["go"],
//   4: ["test"],
//   5: ["hello", "world"],
// }

numbers := []int{1, 2, 3, 4, 5, 6}
evenOdd := slicex.Group(numbers, func(n int) string {
    if n%2 == 0 {
        return "even"
    }
    return "odd"
})
// Result: map[string][]int{
//   "even": [2, 4, 6],
//   "odd":  [1, 3, 5],
// }
```

### MapConcurrent

Creates a concurrent map handler with fluent configuration for high-performance parallel processing. Uses function currying pattern for maximum flexibility.

```go
func MapConcurrent[T, R any](mapFunc func(context.Context, T) (R, error)) *MapConcurrentHandler[T, R]
```

**Configuration Methods:**
- `WithConcurrency(n int)` - Sets maximum concurrent operations (default: 8)
- `WithStopOnError(stop bool)` - Stop on first error (true) or collect all errors (false, default: true)
- `Execute(ctx context.Context, slice []T)` - Runs the concurrent operation

**Example:**
```go
// Basic concurrent mapping
numbers := []int{1, 2, 3, 4, 5}
mapFunc := func(ctx context.Context, n int) (string, error) {
    // Simulate expensive operation
    time.Sleep(100 * time.Millisecond)
    if n == 0 {
        return "", errors.New("cannot process zero")
    }
    return fmt.Sprintf("processed-%d", n*2), nil
}

result, err := slicex.MapConcurrent(mapFunc).
    WithConcurrency(3).
    WithStopOnError(false).
    Execute(context.Background(), numbers)
// Result: ["processed-2", "processed-4", "processed-6", "processed-8", "processed-10"]

// HTTP requests example
urls := []string{"http://api1.com", "http://api2.com", "http://api3.com"}
fetchFunc := func(ctx context.Context, url string) (*http.Response, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    return http.DefaultClient.Do(req)
}

responses, err := slicex.MapConcurrent(fetchFunc).
    WithConcurrency(5).
    WithStopOnError(true).
    Execute(context.Background(), urls)
```

**Key Features:**
- **Order preservation**: Results maintain the same order as input slice
- **Configurable concurrency**: Control maximum parallel operations
- **Error handling strategies**: Stop on first error or collect all errors
- **Context support**: Full context cancellation support
- **Memory efficient**: Uses pre-allocated slices, no mutex needed
- **Fluent API**: Method chaining for clean configuration

## Installation

```bash
go get github.com/letmevibethatforyou/gox/slicex
```

## Usage

```go
import (
    "context"
    "fmt"
    "time"
    "github.com/letmevibethatforyou/gox/slicex"
)

func main() {
    // Example usage
    data := []int{1, 2, 2, 3, 0, 4, 0, 5}
    
    // Get unique non-zero values
    unique := slicex.Unique(slicex.FilterNonZero(data))
    // Result: [1, 2, 3, 4, 5]
    
    // Transform to strings
    strings := slicex.Map(unique, func(n int) string {
        return fmt.Sprintf("num-%d", n)
    })
    // Result: ["num-1", "num-2", "num-3", "num-4", "num-5"]
    
    // Group by string length
    grouped := slicex.Group(strings, func(s string) int {
        return len(s)
    })
    // Result: map[int][]string{
    //   5: ["num-1", "num-2", "num-3", "num-4", "num-5"],
    // }
    
    // Concurrent processing example
    urls := []string{"url1", "url2", "url3", "url4", "url5"}
    processFunc := func(ctx context.Context, url string) (string, error) {
        // Simulate expensive operation
        time.Sleep(100 * time.Millisecond)
        return fmt.Sprintf("processed-%s", url), nil
    }
    
    results, err := slicex.MapConcurrent(processFunc).
        WithConcurrency(3).
        WithStopOnError(false).
        Execute(context.Background(), urls)
    
    if err != nil {
        // Handle error
    }
    // Results processed concurrently while preserving order
}
```

## Requirements

- Go 1.18 or later (for generics support)

## Testing

Run the test suite:

```bash
go test -v ./slicex
```