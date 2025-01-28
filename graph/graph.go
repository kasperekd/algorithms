package graph

import "fmt"

type Edge struct {
	U, V, W int
}

type Graph struct {
	adj   map[int][]int
	edges []Edge
}

func (g *Graph) GetEdges() []Edge {
	edgesCopy := make([]Edge, len(g.edges))
	copy(edgesCopy, g.edges)
	return edgesCopy
}

func NewGraph() *Graph {
	return &Graph{
		adj:   make(map[int][]int),
		edges: make([]Edge, 0),
	}
}

func (g *Graph) AddEdge(u, v, w int) {
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

func (g *Graph) String() string {
	result := ""
	for v, neighbors := range g.adj {
		result += fmt.Sprintf("%d: %v\n", v, neighbors)
	}
	return result
}

func (g *Graph) ConnectedComponents() (count int, comp map[int]int) {
	visited := make(map[int]bool)
	comp = make(map[int]int)
	count = 0

	for v := range g.adj {
		if !visited[v] {
			count++
			dfsOrder := g.DFS(v)

			for _, u := range dfsOrder {
				visited[u] = true
				comp[u] = count
			}

		}
	}

	return count, comp
}

func (g *Graph) GetAllEdges() []Edge {
	var edges []Edge
	for u := range g.adj {
		for _, v := range g.adj[u] {
			if u < v {
				w := 0
				for _, edge := range g.edges {
					if (edge.U == u && edge.V == v) || (edge.U == v && edge.V == u) {
						w = edge.W
						break
					}
				}
				edges = append(edges, Edge{u, v, w})
			}
		}
	}
	return edges
}

func (g *Graph) GetNeighbors(u int) []struct{ v, w int } {
	neighbors := make([]struct{ v, w int }, 0)
	for _, edge := range g.edges {
		if edge.U == u {
			neighbors = append(neighbors, struct{ v, w int }{v: edge.V, w: edge.W})
		}
		if edge.V == u {
			neighbors = append(neighbors, struct{ v, w int }{v: edge.U, w: edge.W})
		}
	}
	return neighbors
}
