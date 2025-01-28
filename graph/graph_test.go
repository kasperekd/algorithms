package graph_test

import (
	"reflect"
	"testing"

	"github.com/kasperekd/algorithms/graph"
)

func TestHasEdge(t *testing.T) {
	g := graph.NewGraph()
	g.AddEdge(0, 1, 0)
	g.AddEdge(0, 2, 0)
	g.AddEdge(1, 2, 0)
	g.AddEdge(2, 3, 0)

	tests := []struct {
		name string
		u    int
		v    int
		want bool
	}{
		{"Existing edge 0-1", 0, 1, true},
		{"Existing edge 1-0", 1, 0, true},
		{"Non-existing edge 0-3", 0, 3, false},
		{"Existing edge 2-3", 2, 3, true},
		{"Existing edge 3-2", 3, 2, true},
		{"Non-existing vertex 4", 4, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := graph.HasEdge(g, tt.u, tt.v); got != tt.want {
				t.Errorf("HasEdge(%d, %d) = %v, want %v", tt.u, tt.v, got, tt.want)
			}
		})
	}

	g2 := graph.NewGraph()
	g2.AddEdge(5, 7, 0)
	if got := graph.HasEdge(g2, 5, 7); !got {
		t.Errorf("HasEdge(5, 7) on g2 = %v, want %v", got, true)
	}
	if got := graph.HasEdge(g2, 5, 8); got {
		t.Errorf("HasEdge(5, 8) on g2 = %v, want %v", got, false)
	}
}

func TestBFS(t *testing.T) {
	g := graph.NewGraph()
	g.AddEdge(0, 1, 0)
	g.AddEdge(0, 2, 0)
	g.AddEdge(1, 2, 0)
	g.AddEdge(1, 3, 0)
	g.AddEdge(2, 4, 0)
	g.AddEdge(3, 5, 0)
	g.AddEdge(5, 6, 0)
	g.AddEdge(6, 4, 0)

	tests := []struct {
		name    string
		start   int
		want    []int
		wantErr bool
	}{
		{"BFS from 0", 0, []int{0, 1, 2, 3, 4, 5, 6}, false},
		{"BFS from 2", 2, []int{2, 0, 1, 4, 3, 6, 5}, false},
		{"BFS from 5", 5, []int{5, 3, 6, 1, 4, 0, 2}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.BFS(tt.start)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BFS(%d) = %v, want %v", tt.start, got, tt.want)
			}
		})
	}

	g2 := graph.NewGraph()
	g2.AddEdge(5, 7, 0)
	if got := g2.BFS(5); !reflect.DeepEqual(got, []int{5, 7}) {
		t.Errorf("BFS(5) on g2 = %v, want %v", got, []int{5, 7})
	}
}

func TestDFS(t *testing.T) {
	g := graph.NewGraph()
	g.AddEdge(0, 1, 0)
	g.AddEdge(0, 2, 0)
	g.AddEdge(1, 3, 0)
	g.AddEdge(1, 4, 0)
	g.AddEdge(2, 5, 0)
	g.AddEdge(3, 5, 0)

	want := []int{0, 1, 3, 5, 2, 4}
	got := g.DFS(0)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DFS(0) = %v, want %v", got, want)
	}
}

func TestConnectedComponents(t *testing.T) {
	g := graph.NewGraph()
	g.AddEdge(0, 1, 0)
	g.AddEdge(0, 2, 0)
	g.AddEdge(1, 3, 0)
	g.AddEdge(4, 5, 0)
	g.AddEdge(5, 6, 0)
	g.AddEdge(6, 4, 0)

	count, comp := g.ConnectedComponents()
	if count != 2 {
		t.Errorf("ConnectedComponents() count = %d, want %d", count, 2)
	}
	expectedComp := map[int]int{0: 1, 1: 1, 2: 1, 3: 1, 4: 2, 5: 2, 6: 2}
	if !reflect.DeepEqual(comp, expectedComp) {
		t.Errorf("ConnectedComponents() comp = %v, want %v", comp, expectedComp)
	}

	g2 := graph.NewGraph()
	g2.AddEdge(0, 1, 0)
	g2.AddEdge(0, 2, 0)
	g2.AddEdge(1, 3, 0)
	g2.AddEdge(2, 4, 0)
	g2.AddEdge(3, 5, 0)
	g2.AddEdge(4, 6, 0)
	g2.AddEdge(5, 6, 0)
	count2, _ := g2.ConnectedComponents()
	if count2 != 1 {
		t.Errorf("ConnectedComponents() count2 = %d, want %d", count2, 1)
	}

	g3 := graph.NewGraph()
	g3.AddEdge(0, 1, 0)
	count3, _ := g3.ConnectedComponents()
	if count3 != 1 {
		t.Errorf("ConnectedComponents() count3 = %d, want %d", count3, 1)
	}
}

func TestGetNeighbors(t *testing.T) {
	g := graph.NewGraph()
	g.AddEdge(0, 1, 4)
	g.AddEdge(0, 2, 2)
	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 5)

	tests := []struct {
		name    string
		vertex  int
		want    []struct{ v, w int }
		wantErr bool
	}{
		{
			name:   "Neighbors of vertex 0",
			vertex: 0,
			want: []struct{ v, w int }{
				{v: 1, w: 4},
				{v: 2, w: 2},
			},
			wantErr: false,
		},
		{
			name:   "Neighbors of vertex 2",
			vertex: 2,
			want: []struct{ v, w int }{
				{v: 0, w: 2},
				{v: 1, w: 1},
				{v: 3, w: 5},
			},
			wantErr: false,
		},
		{
			name:   "Neighbors of vertex 3",
			vertex: 3,
			want: []struct{ v, w int }{
				{v: 2, w: 5},
			},
			wantErr: false,
		},
		{
			name:    "Non-existent vertex",
			vertex:  4,
			want:    []struct{ v, w int }{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.GetNeighbors(tt.vertex)

			if len(got) != len(tt.want) {
				t.Errorf("GetNeighbors() returned a slice of different length than expected. Got: %v, Want: %v", got, tt.want)
				return
			}

			for _, wantItem := range tt.want {
				if !contains(got, wantItem) {
					t.Errorf("GetNeighbors() doesn't contain expected item: %v", wantItem)
				}
			}
		})
	}
}

func contains(s []struct{ V, W int }, item struct{ v, w int }) bool {
	for _, v := range s {
		if struct {
			v int
			w int
		}{v: v.V, w: v.W} == item {
			return true
		}
	}
	return false
}
