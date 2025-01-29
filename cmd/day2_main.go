package main

// import (
// 	"fmt"

// 	"github.com/kasperekd/algorithms/graph"
// 	"github.com/kasperekd/algorithms/unionfind"
// )

// func main() {
// 	ds := unionfind.NewDisjointSet(6)
// 	ds.Union(0, 1)
// 	ds.Union(1, 2)
// 	ds.Union(3, 4)

// 	fmt.Println("DisjointSet:")
// 	for i := 0; i < 6; i++ {
// 		fmt.Printf("Find(%d) = %d ", i, ds.Find(i))
// 	}
// 	fmt.Println()

// 	ds.Union(2, 3)
// 	parents := ds.GetParents()
// 	ranks := ds.GetRanks()

// 	fmt.Println("After Union(2,3):")
// 	fmt.Println("Parents:", parents)
// 	fmt.Println("Ranks:", ranks)

// 	edges := []graph.Edge{
// 		{U: 0, V: 1, W: 4},
// 		{U: 0, V: 2, W: 2},
// 		{U: 1, V: 2, W: 1},
// 		{U: 1, V: 3, W: 5},
// 		{U: 2, V: 3, W: 8},
// 	}

// 	mst, totalWeight := graph.MST(4, edges)
// 	fmt.Println("\nMST:")
// 	fmt.Println("Edges:", mst)
// 	fmt.Println("Total weight:", totalWeight)

// 	g := graph.NewGraph()
// 	g.AddEdge(0, 1, 4)
// 	g.AddEdge(0, 2, 2)
// 	g.AddEdge(1, 2, 1)
// 	g.AddEdge(1, 3, 5)
// 	g.AddEdge(2, 3, 8)

// 	fmt.Println("\nDijkstra:")
// 	distDijkstra, _ := graph.Dijkstra(g, 0)
// 	fmt.Println("Distances:", distDijkstra)

// 	g2 := graph.NewGraph()
// 	g2.AddEdge(0, 1, -1)
// 	g2.AddEdge(0, 2, 4)
// 	g2.AddEdge(1, 2, 3)
// 	g2.AddEdge(1, 3, 2)
// 	g2.AddEdge(3, 1, -5)
// 	g2.AddEdge(2, 0, 4)
// 	g2.AddEdge(3, 0, 2)

// 	fmt.Println("\nBellman-Ford (Negative Weight Graph):")
// 	edges2 := g2.GetEdges()
// 	distBellman2, _, hasNegativeCycle2 := graph.BellmanFord(4, edges2, 0)

// 	if hasNegativeCycle2 {
// 		fmt.Println("Graph has a negative cycle")
// 	} else {
// 		fmt.Println("Distances:", distBellman2)
// 	}
// }
