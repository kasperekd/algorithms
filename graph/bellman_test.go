package graph

import (
	"reflect"
	"testing"
)

func TestBellmanFord(t *testing.T) {
	tests := []struct {
		name              string
		n                 int // Количество вершин
		edges             []Edge
		start             int
		wantDist          []int
		wantParent        []int
		wantNegativeCycle bool
	}{
		{
			name: "Graph with negative cycle",
			n:    4,
			edges: []Edge{
				{U: 0, V: 1, W: -1},
				{U: 0, V: 2, W: 4},
				{U: 1, V: 2, W: 3},
				{U: 1, V: 3, W: 2},
				{U: 3, V: 1, W: -5},
			},
			start:             0,
			wantDist:          nil,
			wantParent:        nil,
			wantNegativeCycle: true,
		},
		{
			name: "Disconnected graph",
			n:    6,
			edges: []Edge{
				{U: 0, V: 1, W: 1},
				{U: 0, V: 2, W: 2},
				{U: 3, V: 4, W: 1},
				{U: 3, V: 5, W: 2},
			},
			start:             0,
			wantDist:          []int{0, 1, 2, 2147483647, 2147483647, 2147483647},
			wantParent:        []int{-1, 0, 0, -1, -1, -1},
			wantNegativeCycle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist, parent, hasNegativeCycle := BellmanFord(tt.n, tt.edges, tt.start)

			if !reflect.DeepEqual(dist, tt.wantDist) {
				t.Errorf("BellmanFord() dist = %+v, want %+v", dist, tt.wantDist)
			}
			if !reflect.DeepEqual(parent, tt.wantParent) {
				t.Errorf("BellmanFord() parent = %+v, want %+v", parent, tt.wantParent)
			}
			if hasNegativeCycle != tt.wantNegativeCycle {
				t.Errorf("BellmanFord() hasNegativeCycle = %v, want %v", hasNegativeCycle, tt.wantNegativeCycle)
			}
		})
	}
}
