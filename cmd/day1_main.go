package main

import (
	"fmt"

	"github.com/kasperekd/algorithms/graph"
)

func task3() {
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

func task4() {
	g := graph.NewGraph()

	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 5)
	g.AddEdge(5, 6)
	g.AddEdge(6, 4)

	g.PrintGraph()

	bfsOrder := g.BFS(0)
	fmt.Println("BFS начиная с 0:", bfsOrder)

	bfsOrder2 := g.BFS(2)
	fmt.Println("BFS начиная с 2:", bfsOrder2)

	bfsOrder3 := g.BFS(5)
	fmt.Println("BFS начиная с 5:", bfsOrder3)

	g2 := graph.NewGraph()
	g2.AddEdge(5, 7)
	bfsOrder4 := g2.BFS(5)
	fmt.Println("BFS g2 начиная с 5:", bfsOrder4)

	g3 := graph.NewGraph()
	bfsOrder5 := g3.BFS(0)
	fmt.Println("BFS обход пустого графа g3 с 0:", bfsOrder5)
}

func main() {
	task3()
	task4()
}
