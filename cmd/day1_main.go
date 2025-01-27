package main

func main() {
	g := NewGraph()

	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	g.PrintGraph()

	g2 := NewGraph()
	g2.AddEdge(5, 7)
	g2.PrintGraph()
}
