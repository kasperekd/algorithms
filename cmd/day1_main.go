package main

import (
	"fmt"

	"github.com/kasperekd/algorithms/graph"
)

func main() {
	g := graph.NewGraph()

	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	g.PrintGraph()

	// HasEdge
	fmt.Println("HasEdge(g, 0, 1):", graph.HasEdge(g, 0, 1))
	fmt.Println("HasEdge(g, 0, 1):", graph.HasEdge(g, 0, 1))
	fmt.Println("HasEdge(g, 1, 0):", graph.HasEdge(g, 1, 0))
	fmt.Println("HasEdge(g, 0, 3):", graph.HasEdge(g, 0, 3))
	fmt.Println("HasEdge(g, 2, 3):", graph.HasEdge(g, 2, 3))
	fmt.Println("HasEdge(g, 3, 2):", graph.HasEdge(g, 3, 2))
	fmt.Println("HasEdge(g, 4, 0):", graph.HasEdge(g, 4, 0))

	g2 := graph.NewGraph()
	g2.AddEdge(5, 7)
	g2.PrintGraph()
	fmt.Println("HasEdge(g2, 5, 7):", graph.HasEdge(g2, 5, 7))
	fmt.Println("HasEdge(g2, 7, 5):", graph.HasEdge(g2, 7, 5))
	fmt.Println("HasEdge(g2, 5, 8):", graph.HasEdge(g2, 5, 8))
}
