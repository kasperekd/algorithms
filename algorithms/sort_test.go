package algorithms

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestMergeSort(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "Empty array",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "Sorted array",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Unsorted array",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Array with duplicates",
			input:    []int{5, 4, 3, 2, 1, 5, 4, 3, 2, 1},
			expected: []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5},
		},
		{
			name:     "Large array",
			input:    generateRandomArray(100000),
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := parallelMergeSort(tc.input)

			if tc.name == "Large array" {
				expected := make([]int, len(tc.input))
				copy(expected, tc.input)
				sort.Ints(expected)

				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("Expected sorted array, but got %v", actual)
				}
			} else {
				if !reflect.DeepEqual(actual, tc.expected) {
					t.Errorf("Expected %v, but got %v", tc.expected, actual)
				}
			}
		})
	}
}

func generateRandomArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(size * 10)
	}
	return arr
}
