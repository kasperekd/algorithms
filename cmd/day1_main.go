package main

// func main() {
// 	g := graph.NewGraph()

// 	g.AddEdge(0, 1, 0)
// 	g.AddEdge(0, 2, 0)
// 	g.AddEdge(1, 3, 0)
// 	g.AddEdge(2, 4, 0)
// 	g.AddEdge(3, 5, 0)

// 	fmt.Println("Структура графа g:")
// 	fmt.Println(g.String())

// 	fmt.Println("HasEdge(g, 0, 1):", graph.HasEdge(g, 0, 1))
// 	fmt.Println("HasEdge(g, 1, 0):", graph.HasEdge(g, 1, 0))
// 	fmt.Println("HasEdge(g, 0, 3):", graph.HasEdge(g, 0, 3))
// 	fmt.Println("HasEdge(g, 3, 0):", graph.HasEdge(g, 3, 0))

// 	bfsOrder := g.BFS(0)
// 	fmt.Println("BFS(g, 0):", bfsOrder)

// 	dfsOrder := g.DFS(0)
// 	fmt.Println("DFS(g, 0):", dfsOrder)

// 	// Острова
// 	g2 := graph.NewGraph()
// 	g2.AddEdge(0, 1, 0)
// 	g2.AddEdge(2, 3, 0)
// 	g2.AddEdge(4, 5, 0)
// 	g2.AddEdge(5, 6, 0)

// 	fmt.Println("\ng2 острова:")
// 	fmt.Println(g2.String())

// 	// ConnectedComponents
// 	count, comp := g2.ConnectedComponents()
// 	fmt.Println("Связных компонент g2:", count)
// 	fmt.Println("Сопоставление с g2:", comp)

// 	g3 := graph.NewGraph()
// 	g3.AddEdge(0, 1, 0)
// 	fmt.Println("\nпростой:")
// 	fmt.Println(g3.String())
// 	count3, comp3 := g3.ConnectedComponents()
// 	fmt.Println("Связных компонент g3:", count3)
// 	fmt.Println("Сопоставление g3:", comp3)
// }
