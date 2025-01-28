package graph

func BellmanFord(n int, edges []Edge, start int) ([]int, []int, bool) {
	dist := make([]int, n)
	parent := make([]int, n)

	for i := 0; i < n; i++ {
		dist[i] = 1<<31 - 1
		parent[i] = -1
	}
	dist[start] = 0

	for i := 1; i < n; i++ {
		for _, edge := range edges {
			if dist[edge.U] != 1<<31-1 && dist[edge.U]+edge.W < dist[edge.V] {
				dist[edge.V] = dist[edge.U] + edge.W
				parent[edge.V] = edge.U
			}
		}
	}

	for _, edge := range edges {
		if dist[edge.U] != 1<<31-1 && dist[edge.U]+edge.W < dist[edge.V] {
			return nil, nil, true
		}
	}

	return dist, parent, false
}
