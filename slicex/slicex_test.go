package slicex

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestUnique(t *testing.T) {
	tests := map[string]struct {
		input    []int
		expected []int
	}{
		"empty slice": {
			input:    []int{},
			expected: nil,
		},
		"single element": {
			input:    []int{1},
			expected: []int{1},
		},
		"all unique": {
			input:    []int{1, 2, 3, 4},
			expected: []int{1, 2, 3, 4},
		},
		"with duplicates": {
			input:    []int{1, 2, 2, 3, 1, 4, 3},
			expected: []int{1, 2, 3, 4},
		},
		"all same": {
			input:    []int{5, 5, 5, 5},
			expected: []int{5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Unique(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Unique(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUniqueStrings(t *testing.T) {
	input := []string{"hello", "world", "hello", "go", "world", "test"}
	expected := []string{"hello", "world", "go", "test"}
	result := Unique(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unique(%v) = %v, expected %v", input, result, expected)
	}
}

func TestFilterNonZero(t *testing.T) {
	tests := map[string]struct {
		input    []int
		expected []int
	}{
		"empty slice": {
			input:    []int{},
			expected: []int{},
		},
		"no zeros": {
			input:    []int{1, 2, 3, 4},
			expected: []int{1, 2, 3, 4},
		},
		"with zeros": {
			input:    []int{1, 0, 2, 0, 3, 0, 4},
			expected: []int{1, 2, 3, 4},
		},
		"all zeros": {
			input:    []int{0, 0, 0, 0},
			expected: []int{},
		},
		"negative numbers": {
			input:    []int{-1, 0, -2, 0, 3},
			expected: []int{-1, -2, 3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := FilterNonZero(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FilterNonZero(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFilterNonZeroStrings(t *testing.T) {
	input := []string{"hello", "", "world", "", "test"}
	expected := []string{"hello", "world", "test"}
	result := FilterNonZero(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FilterNonZero(%v) = %v, expected %v", input, result, expected)
	}
}

func TestMap(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		expected := []string{"1", "2", "3", "4"}
		result := Map(input, func(i int) string {
			return strconv.Itoa(i)
		})

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Map(%v, intToString) = %v, expected %v", input, result, expected)
		}
	})

	t.Run("string to length", func(t *testing.T) {
		input := []string{"hello", "world", "go", "test"}
		expected := []int{5, 5, 2, 4}
		result := Map(input, func(s string) int {
			return len(s)
		})

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Map(%v, stringToLength) = %v, expected %v", input, result, expected)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		expected := []string(nil)
		result := Map(input, func(i int) string {
			return strconv.Itoa(i)
		})

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Map(%v, intToString) = %v, expected %v", input, result, expected)
		}
	})

	t.Run("square numbers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []int{1, 4, 9, 16, 25}
		result := Map(input, func(i int) int {
			return i * i
		})

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Map(%v, square) = %v, expected %v", input, result, expected)
		}
	})
}

func TestGroup(t *testing.T) {
	t.Run("group by string length", func(t *testing.T) {
		input := []string{"hello", "world", "go", "test", "a", "b"}
		result := Group(input, func(s string) int {
			return len(s)
		})

		expected := map[int][]string{
			1: {"a", "b"},
			2: {"go"},
			4: {"test"},
			5: {"hello", "world"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Group(%v, lengthKey) = %v, expected %v", input, result, expected)
		}
	})

	t.Run("group by even/odd", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		result := Group(input, func(i int) string {
			if i%2 == 0 {
				return "even"
			}
			return "odd"
		})

		expected := map[string][]int{
			"even": {2, 4, 6},
			"odd":  {1, 3, 5},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Group(%v, evenOddKey) = %v, expected %v", input, result, expected)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		result := Group(input, func(i int) string {
			return "key"
		})

		expected := map[string][]int{}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Group(%v, constantKey) = %v, expected %v", input, result, expected)
		}
	})

	t.Run("single group", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		result := Group(input, func(i int) string {
			return "same"
		})

		expected := map[string][]int{
			"same": {1, 2, 3, 4},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Group(%v, constantKey) = %v, expected %v", input, result, expected)
		}
	})
}

type Person struct {
	Name string
	Age  int
}

func TestGroupComplex(t *testing.T) {
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 30},
		{"Diana", 25},
		{"Eve", 35},
	}

	result := Group(people, func(p Person) int {
		return p.Age
	})

	expected := map[int][]Person{
		25: {{"Bob", 25}, {"Diana", 25}},
		30: {{"Alice", 30}, {"Charlie", 30}},
		35: {{"Eve", 35}},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Group people by age failed: got %v, expected %v", result, expected)
	}
}

func TestMapConcurrent(t *testing.T) {
	t.Run("basic concurrent execution", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}

		mapFunc := func(ctx context.Context, n int) (string, error) {
			time.Sleep(10 * time.Millisecond) // Simulate work
			return strconv.Itoa(n * 2), nil
		}

		vv, err := MapConcurrent(mapFunc).Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		expected := []string{"2", "4", "6", "8", "10"}
		if !reflect.DeepEqual(vv, expected) {
			t.Errorf("Expected %v, got %v", expected, vv)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}

		mapFunc := func(ctx context.Context, n int) (string, error) {
			return strconv.Itoa(n), nil
		}

		result, err := MapConcurrent(mapFunc).Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil mapConcurrentResult for empty input, got %v", result)
		}
	})

	t.Run("with custom concurrency", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		mapFunc := func(ctx context.Context, n int) (int, error) {
			time.Sleep(10 * time.Millisecond) // Simulate work
			return n * n, nil
		}

		result, err := MapConcurrent(mapFunc).
			WithConcurrency(3).
			Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		expected := []int{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("stop on first error", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}

		mapFunc := func(ctx context.Context, n int) (int, error) {
			if n == 3 {
				return 0, errors.New("error at 3")
			}
			time.Sleep(50 * time.Millisecond) // Simulate work
			return n * 2, nil
		}

		result, err := MapConcurrent(mapFunc).
			WithStopOnError(true).
			Execute(context.Background(), input)

		if err == nil {
			t.Fatal("Expected error but got none")
		}

		if err.Error() != "error at 3" {
			t.Errorf("Expected 'error at 3', got '%v'", err)
		}

		// Result should be nil when there's an error
		if result != nil {
			t.Errorf("Expected nil mapConcurrentResult when error occurs, got %v", result)
		}
	})

	t.Run("continue on error", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}

		mapFunc := func(ctx context.Context, n int) (int, error) {
			if n == 3 || n == 4 {
				return 0, errors.New("error at " + strconv.Itoa(n))
			}
			return n * 2, nil
		}

		result, err := MapConcurrent(mapFunc).
			WithStopOnError(false).
			Execute(context.Background(), input)

		if err == nil {
			t.Fatal("Expected error but got none")
		}

		// Result should be nil when there's an error
		if result != nil {
			t.Errorf("Expected nil mapConcurrentResult when error occurs, got %v", result)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}

		mapFunc := func(ctx context.Context, n int) (int, error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return n * 2, nil
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		result, err := MapConcurrent(mapFunc).Execute(ctx, input)

		if err == nil {
			t.Fatal("Expected context cancellation error but got none")
		}

		// Result should be nil when there's an error (including cancellation)
		if result != nil {
			t.Errorf("Expected nil mapConcurrentResult when context is cancelled, got %v", result)
		}
	})

	t.Run("order preservation", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		mapFunc := func(ctx context.Context, n int) (string, error) {
			// Add variable delay to test ordering
			delay := time.Duration((11-n)*10) * time.Millisecond
			time.Sleep(delay)
			return "item-" + strconv.Itoa(n), nil
		}

		result, err := MapConcurrent(mapFunc).
			WithConcurrency(5).
			Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		expected := []string{
			"item-1", "item-2", "item-3", "item-4", "item-5",
			"item-6", "item-7", "item-8", "item-9", "item-10",
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Order not preserved. Expected %v, got %v", expected, result)
		}
	})
}

func TestMapConcurrentWorkerPool(t *testing.T) {
	t.Run("worker pool uses minimum of concurrency and slice length", func(t *testing.T) {
		// Test with slice smaller than concurrency
		smallInput := []int{1, 2}
		concurrentCount := 0
		maxConcurrent := 0
		var mu sync.Mutex

		mapFunc := func(ctx context.Context, n int) (int, error) {
			mu.Lock()
			concurrentCount++
			if concurrentCount > maxConcurrent {
				maxConcurrent = concurrentCount
			}
			mu.Unlock()

			time.Sleep(50 * time.Millisecond) // Hold the worker for a bit

			mu.Lock()
			concurrentCount--
			mu.Unlock()

			return n * 2, nil
		}

		result, err := MapConcurrent(mapFunc).
			WithConcurrency(8). // Much higher than slice length
			Execute(context.Background(), smallInput)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		expected := []int{2, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Should only use 2 workers max (slice length), not 8
		if maxConcurrent > len(smallInput) {
			t.Errorf("Expected max %d concurrent workers, but saw %d", len(smallInput), maxConcurrent)
		}

		// Should have used exactly the slice length as worker count
		if maxConcurrent != len(smallInput) {
			t.Errorf("Expected exactly %d concurrent workers, but saw %d", len(smallInput), maxConcurrent)
		}
	})

	t.Run("worker pool respects concurrency limit for large slices", func(t *testing.T) {
		// Test with slice larger than concurrency
		largeInput := make([]int, 20)
		for i := range largeInput {
			largeInput[i] = i + 1
		}

		concurrentCount := 0
		maxConcurrent := 0
		var mu sync.Mutex

		mapFunc := func(ctx context.Context, n int) (int, error) {
			mu.Lock()
			concurrentCount++
			if concurrentCount > maxConcurrent {
				maxConcurrent = concurrentCount
			}
			mu.Unlock()

			time.Sleep(20 * time.Millisecond) // Hold the worker for a bit

			mu.Lock()
			concurrentCount--
			mu.Unlock()

			return n * 2, nil
		}

		concurrencyLimit := 3
		result, err := MapConcurrent(mapFunc).
			WithConcurrency(concurrencyLimit).
			Execute(context.Background(), largeInput)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify all results are correct
		if len(result) != len(largeInput) {
			t.Errorf("Expected mapConcurrentResult length %d, got %d", len(largeInput), len(result))
		}

		// Should never exceed concurrency limit
		if maxConcurrent > concurrencyLimit {
			t.Errorf("Expected max %d concurrent workers, but saw %d", concurrencyLimit, maxConcurrent)
		}

		// Should have used exactly the concurrency limit
		if maxConcurrent != concurrencyLimit {
			t.Errorf("Expected exactly %d concurrent workers, but saw %d", concurrencyLimit, maxConcurrent)
		}
	})
}

func TestMapConcurrentFluentAPI(t *testing.T) {
	t.Run("method chaining", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}

		mapFunc := func(ctx context.Context, n int) (int, error) {
			return n * 3, nil
		}

		result, err := MapConcurrent(mapFunc).
			WithConcurrency(2).
			WithStopOnError(false).
			Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		expected := []int{3, 6, 9, 12, 15}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("default configuration", func(t *testing.T) {
		mapFunc := func(ctx context.Context, n int) (int, error) {
			return n + 10, nil
		}

		handler := MapConcurrent(mapFunc)

		// Check defaults
		if handler.concurrency != 8 {
			t.Errorf("Expected default concurrency 8, got %d", handler.concurrency)
		}

		if !handler.stopOnError {
			t.Error("Expected default stopOnError to be true")
		}
	})
}
