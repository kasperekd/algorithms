package graph

func (g *Graph) BFS(start int) []int {
	visited := make(map[int]bool)
	var queue SimpleQueue
	order := []int{}

	queue.Enqueue(start)
	visited[start] = true

	for len(queue) > 0 {
		u, ok := queue.Dequeue()
		if !ok {
			break
		}

		order = append(order, u)

		if neighbors, ok := g.adj[u]; ok {
			for _, v := range neighbors {
				if !visited[v] {
					visited[v] = true
					queue.Enqueue(v)
				}
			}
		}
	}

	return order
}
