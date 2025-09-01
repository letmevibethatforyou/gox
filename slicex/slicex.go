// Package slicex provides utility functions for slice operations with generics support.
package slicex

// Unique returns a new slice containing only unique elements from the input slice,
// preserving the order of first occurrence.
func Unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}

	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// FilterNonZero returns a new slice with all non-zero values from the input slice.
// Zero values are determined by Go's zero value concept (0, "", nil, etc.).
func FilterNonZero[T comparable](slice []T) []T {
	var zero T
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if item != zero {
			result = append(result, item)
		}
	}

	return result
}

// Map applies the given function to each element of the slice and returns
// a new slice containing the results.
func Map[T, R any](slice []T, fn func(T) R) []R {
	if len(slice) == 0 {
		return nil
	}

	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}

	return result
}

// Group groups the elements of the slice by the result of the key function.
// Returns a map where keys are the grouping criteria and values are slices
// of grouped items.
func Group[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)

	for _, item := range slice {
		key := keyFn(item)
		result[key] = append(result[key], item)
	}

	return result
}
