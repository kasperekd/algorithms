package graph

import (
	"reflect"
	"sort"
	"testing"
)

func TestMST(t *testing.T) {
	tests := []struct {
		name       string
		n          int
		edges      []Edge
		wantMST    []Edge
		wantWeight int
	}{
		{
			name: "Simple graph",
			n:    4,
			edges: []Edge{
				{U: 0, V: 1, W: 4},
				{U: 0, V: 2, W: 2},
				{U: 1, V: 2, W: 1},
				{U: 1, V: 3, W: 5},
				{U: 2, V: 3, W: 8},
			},
			wantMST: []Edge{
				{U: 1, V: 2, W: 1},
				{U: 0, V: 2, W: 2},
				{U: 1, V: 3, W: 5},
			},
			wantWeight: 8,
		},
		{
			name: "Graph with equal weights",
			n:    4,
			edges: []Edge{
				{U: 0, V: 1, W: 1},
				{U: 0, V: 2, W: 1},
				{U: 1, V: 2, W: 1},
				{U: 1, V: 3, W: 1},
				{U: 2, V: 3, W: 1},
			},
			wantMST: []Edge{
				{U: 0, V: 1, W: 1},
				{U: 0, V: 2, W: 1},
				{U: 1, V: 3, W: 1},
			},
			wantWeight: 3,
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
			wantMST: []Edge{
				{U: 0, V: 1, W: 1},
				{U: 0, V: 2, W: 2},
				{U: 3, V: 4, W: 1},
				{U: 3, V: 5, W: 2},
			},
			wantWeight: 6,
		},
		{
			name: "Complete graph",
			n:    5,
			edges: []Edge{
				{U: 0, V: 1, W: 1},
				{U: 0, V: 2, W: 2},
				{U: 0, V: 3, W: 3},
				{U: 0, V: 4, W: 4},
				{U: 1, V: 2, W: 5},
				{U: 1, V: 3, W: 6},
				{U: 1, V: 4, W: 7},
				{U: 2, V: 3, W: 8},
				{U: 2, V: 4, W: 9},
				{U: 3, V: 4, W: 10},
			},
			wantMST: []Edge{
				{U: 0, V: 1, W: 1},
				{U: 0, V: 2, W: 2},
				{U: 0, V: 3, W: 3},
				{U: 0, V: 4, W: 4},
			},
			wantWeight: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMST, gotWeight := MST(tt.n, tt.edges)

			sortEdges(gotMST)
			sortEdges(tt.wantMST)

			if !reflect.DeepEqual(gotMST, tt.wantMST) {
				t.Errorf("MST() gotMST = %v, want %v", gotMST, tt.wantMST)
			}
			if gotWeight != tt.wantWeight {
				t.Errorf("MST() gotWeight = %v, want %v", gotWeight, tt.wantWeight)
			}
		})
	}
}

func sortEdges(edges []Edge) {
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].U != edges[j].U {
			return edges[i].U < edges[j].U
		}
		return edges[i].V < edges[j].V
	})
}
