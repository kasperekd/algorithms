package algorithms

import (
	"sync"
)

type Edge struct {
	U, V, W int
}

func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := arr[:mid]
	right := arr[mid:]

	left = mergeSort(left)
	right = mergeSort(right)

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

func parallelMergeSort(arr []int) []int {
	if len(arr) <= 1000 {
		return mergeSort(arr)
	}

	mid := len(arr) / 2
	left := arr[:mid]
	right := arr[mid:]

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		left = parallelMergeSort(left)
	}()

	go func() {
		defer wg.Done()
		right = parallelMergeSort(right)
	}()

	wg.Wait()

	return merge(left, right)
}

func ParallelMergeSortEdges(edges []Edge) []Edge {
	if len(edges) <= 1 {
		return edges
	}

	mid := len(edges) / 2
	var left []Edge
	var right []Edge

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		left = ParallelMergeSortEdges(edges[:mid])
	}()

	go func() {
		defer wg.Done()
		right = ParallelMergeSortEdges(edges[mid:])
	}()

	wg.Wait()

	return mergeEdges(left, right)
}

func mergeEdges(left, right []Edge) []Edge {
	result := make([]Edge, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i].W <= right[j].W {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}
