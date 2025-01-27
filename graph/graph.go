package graph

import "fmt"

type Graph struct {
	adj map[int][]int
}

func NewGraph() *Graph {
	return &Graph{
		adj: make(map[int][]int),
	}
}

func (g *Graph) AddEdge(u, v int) {
	if _, ok := g.adj[u]; !ok {
		g.adj[u] = make([]int, 0)
	}
	g.adj[u] = append(g.adj[u], v)

	if _, ok := g.adj[v]; !ok {
		g.adj[v] = make([]int, 0)
	}
	g.adj[v] = append(g.adj[v], u)
}

func (g *Graph) PrintGraph() {
	for vertex, neighbors := range g.adj {
		fmt.Printf("Vertex %d: Neighbors: %v\n", vertex, neighbors)
	}
}

func HasEdge(g *Graph, u, v int) bool {
	if neighbors, ok := g.adj[u]; ok {
		for _, neighbor := range neighbors {
			if neighbor == v {
				return true
			}
		}
	}
	return false
}
