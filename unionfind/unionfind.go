package unionfind

type DisjointSet struct {
	parent []int
	rank   []int
	// или size []int — по желанию
}

func (ds *DisjointSet) GetParents() []int {
	parentsCopy := make([]int, len(ds.parent))
	copy(parentsCopy, ds.parent)
	return parentsCopy
}

func (ds *DisjointSet) GetRanks() []int {
	ranksCopy := make([]int, len(ds.rank))
	return ranksCopy
}

func NewDisjointSet(n int) *DisjointSet {
	d := &DisjointSet{
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.rank[i] = 0
	}
	return d
}

func (d *DisjointSet) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DisjointSet) Union(x, y int) bool {
	rx := d.Find(x)
	ry := d.Find(y)
	if rx == ry {
		return false
	}

	// by rank
	if d.rank[rx] < d.rank[ry] {
		d.parent[rx] = ry
	} else if d.rank[rx] > d.rank[ry] {
		d.parent[ry] = rx
	} else {
		d.parent[ry] = rx
		d.rank[rx]++
	}
	return true
}
