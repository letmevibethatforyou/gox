package slicex

import (
	"reflect"
	"strconv"
	"testing"
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
