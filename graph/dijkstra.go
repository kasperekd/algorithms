package graph

import (
	"math"
)

func Dijkstra(g *Graph, start int) ([]int, []int) {
	dist := make([]int, len(g.adj))
	parent := make([]int, len(g.adj))
	visited := make([]bool, len(g.adj))

	vertexMap := make(map[int]int)
	indexMap := make(map[int]int)
	i := 0
	for v := range g.adj {
		vertexMap[v] = i
		indexMap[i] = v
		i++
	}

	for i := 0; i < len(g.adj); i++ {
		dist[i] = math.MaxInt
		parent[i] = -1
		visited[i] = false
	}

	startIndex := vertexMap[start]
	dist[startIndex] = 0

	for count := 0; count < len(g.adj)-1; count++ {
		uIndex := -1
		minDist := math.MaxInt

		for vIndex := 0; vIndex < len(g.adj); vIndex++ {
			if !visited[vIndex] && dist[vIndex] <= minDist {
				minDist = dist[vIndex]
				uIndex = vIndex
			}
		}

		if uIndex == -1 {
			break
		}

		u := indexMap[uIndex]

		visited[uIndex] = true

		for _, neighbor := range g.GetNeighbors(u) {
			v := neighbor.V
			w := neighbor.W
			vIndex := vertexMap[v]

			if !visited[vIndex] && dist[uIndex] != math.MaxInt && dist[uIndex]+w < dist[vIndex] {
				dist[vIndex] = dist[uIndex] + w
				parent[vIndex] = u
			}
		}
	}

	originalParents := make([]int, len(parent))
	for i := range parent {
		if parent[i] != -1 {
			originalParents[i] = indexMap[parent[i]]
		} else {
			originalParents[i] = -1
		}
	}

	return dist, originalParents
}
