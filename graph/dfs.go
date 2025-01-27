package graph

func (g *Graph) DFS(start int) []int {
	visited := make(map[int]bool)
	stack := NewStack()
	order := []int{}

	stack.Push(start)

	for !stack.IsEmpty() {
		u, ok := stack.Pop()
		if !ok {
			continue
		}

		if !visited[u] {
			visited[u] = true
			order = append(order, u)

			if neighbors, ok := g.adj[u]; ok {
				for i := len(neighbors) - 1; i >= 0; i-- {
					v := neighbors[i]
					if !visited[v] {
						stack.Push(v)
					}
				}
			}
		}
	}

	return order
}
