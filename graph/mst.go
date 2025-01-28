package graph

import (
	"github.com/kasperekd/algorithms/algorithms"
	"github.com/kasperekd/algorithms/unionfind"
)

func MST(n int, edges []Edge) (mst []Edge, totalWeight int) {
	sortEdges := make([]algorithms.Edge, len(edges))
	for i, edge := range edges {
		sortEdges[i] = algorithms.Edge{U: edge.U, V: edge.V, W: edge.W}
	}

	sortEdges = algorithms.ParallelMergeSortEdges(sortEdges)

	edges = make([]Edge, len(sortEdges))
	for i, edge := range sortEdges {
		edges[i] = Edge{U: edge.U, V: edge.V, W: edge.W}
	}

	ds := unionfind.NewDisjointSet(n)

	for _, edge := range edges {
		if ds.Find(edge.U) != ds.Find(edge.V) {
			ds.Union(edge.U, edge.V)
			mst = append(mst, edge)
			totalWeight += edge.W
		}
		if len(mst) == n-1 {
			break
		}
	}

	return mst, totalWeight
}
