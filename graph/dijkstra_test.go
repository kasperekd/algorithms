package graph

import (
	"math"
	"reflect"
	"testing"
)

func TestDijkstra(t *testing.T) {
	g1 := NewGraph()
	g1.AddEdge(0, 1, 4)
	g1.AddEdge(0, 7, 8)
	g1.AddEdge(1, 2, 8)
	g1.AddEdge(1, 7, 11)
	g1.AddEdge(2, 3, 7)
	g1.AddEdge(2, 5, 4)
	g1.AddEdge(2, 8, 2)
	g1.AddEdge(3, 4, 9)
	g1.AddEdge(3, 5, 14)
	g1.AddEdge(4, 5, 10)
	g1.AddEdge(5, 6, 2)
	g1.AddEdge(6, 7, 1)
	g1.AddEdge(6, 8, 6)
	g1.AddEdge(7, 8, 7)

	expectedDist1 := []int{0, 4, 12, 19, 28, 16, 18, 8, 10}
	expectedParent1 := []int{-1, 0, 1, 2, 3, 2, 5, 0, 2}
	dist1, parent1 := Dijkstra(g1, 0)

	if !reflect.DeepEqual(dist1, expectedDist1) {
		t.Errorf("Test Case 1 Failed: Expected distances %v, got %v", expectedDist1, dist1)
	}
	if !reflect.DeepEqual(parent1, expectedParent1) {
		t.Errorf("Test Case 1 Failed: Expected parents %v, got %v", expectedParent1, parent1)
	}

	g2 := NewGraph()
	g2.AddEdge(0, 1, 1)
	g2.AddEdge(2, 3, 1)

	expectedDist2 := []int{0, 1, math.MaxInt, math.MaxInt}
	expectedParent2 := []int{-1, 0, -1, -1}
	dist2, parent2 := Dijkstra(g2, 0)

	if !reflect.DeepEqual(dist2, expectedDist2) {
		t.Errorf("Test Case 2 Failed: Expected distances %v, got %v", expectedDist2, dist2)
	}

	if !reflect.DeepEqual(parent2, expectedParent2) {
		t.Errorf("Test Case 2 Failed: Expected parents %v, got %v", expectedParent2, parent2)
	}

	// g3 := NewGraph()
	// g3.AddEdge(0, 1, -1)

	// dist3, _ := Dijkstra(g3, 0)

	// if !reflect.DeepEqual(dist3, []int{0, -1}) {
	// 	t.Logf("Test Case 3: Dijkstra's algorithm doesn't work with negative weights. Got distances %v. Result is undefined for graphs with negative weights", dist3)
	// }

	g4 := NewGraph()
	g4.AddEdge(0, 0, 0)

	expectedDist4 := []int{0}
	expectedParent4 := []int{-1}

	dist4, parent4 := Dijkstra(g4, 0)

	if !reflect.DeepEqual(dist4, expectedDist4) {
		t.Errorf("Test Case 4 Failed: Expected distances %v, got %v", expectedDist4, dist4)
	}

	if !reflect.DeepEqual(parent4, expectedParent4) {
		t.Errorf("Test Case 4 Failed: Expected parents %v, got %v", expectedParent4, parent4)
	}
}
