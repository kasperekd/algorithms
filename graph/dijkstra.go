package graph

import "container/heap"

type Item struct {
	vertex, dist int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func Dijkstra(g *Graph, start int) ([]int, []int) {
	dist := make([]int, len(g.adj))
	parent := make([]int, len(g.adj))
	visited := make([]bool, len(g.adj))

	for i := range dist {
		dist[i] = 1<<31 - 1
		parent[i] = -1
	}

	dist[start] = 0

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{vertex: start, dist: 0})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		u := item.vertex

		if visited[u] {
			continue
		}

		visited[u] = true

		for _, neighbor := range g.GetNeighbors(u) {
			w := neighbor.v
			weight := neighbor.w
			if dist[u]+weight < dist[w] {
				dist[w] = dist[u] + weight
				parent[w] = u
				heap.Push(&pq, &Item{vertex: w, dist: dist[w]})
			}
		}
	}

	return dist, parent
}
