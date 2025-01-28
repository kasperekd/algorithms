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

	found := false
	for _, edge := range g.edges {
		if (edge.U == u && edge.V == v) || (edge.U == v && edge.V == u) {
			found = true
			break
		}
	}
	if !found {
		g.edges = append(g.edges, Edge{U: u, V: v, W: w})
	}

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

func (g *Graph) GetNeighbors(u int) []struct{ V, W int } {
	neighbors := make([]struct{ V, W int }, 0)
	for _, v := range g.adj[u] {
		w := -1
		for _, edge := range g.edges {
			if (edge.U == u && edge.V == v) || (edge.U == v && edge.V == u) {
				w = edge.W
				break
			}
		}
		if w != -1 {
			neighbors = append(neighbors, struct{ V, W int }{V: v, W: w})
		} else {
			fmt.Printf("Warning: Edge (%d, %d) not found in edges list.\n", u, v)
		}
	}
	return neighbors
}
