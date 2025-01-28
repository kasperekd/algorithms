package unionfind

import (
	"testing"
)

func TestDisjointSet(t *testing.T) {
	n := 10
	d := NewDisjointSet(n)

	for i := 0; i < n; i++ {
		if d.Find(i) != i {
			t.Errorf("Find(%d) != %d", i, i)
		}
	}

	d.Union(0, 1)
	d.Union(2, 3)
	d.Union(1, 2)

	if d.Find(0) != d.Find(1) || d.Find(1) != d.Find(2) || d.Find(2) != d.Find(3) {
		t.Error("0, 1, 2, 3 are not in the same component")
	}

	for i := 4; i < n; i++ {
		if d.Find(i) != i {
			t.Errorf("Find(%d) != %d", i, i)
		}
	}

	if d.rank[d.Find(0)] != 2 {
		t.Errorf("Rank != 2")
	}

	if !d.Union(0, 2) {
		t.Log("Union(0, 2) returned false")
	} else {
		t.Error("Union(0, 2) should have returned false")
	}

	largeN := 1000
	largeD := NewDisjointSet(largeN)
	for i := 0; i < largeN-1; i++ {
		largeD.Union(i, i+1)
	}
	if largeD.Find(0) != largeD.Find(largeN-1) {
		t.Error("All elements are not in the same componet")
	}

	emptyD := NewDisjointSet(0)
	if emptyD.parent == nil || emptyD.rank == nil {
		t.Error("parent or rank not initialized for empty set")
	}
}
