package distributed

import (
	"fmt"
	"testing"
	"time"
)

func TestBullyElection(t *testing.T) {
	network := NewNetwork()
	numNodes := 5
	nodes := make([]*Node, numNodes)

	for i := 0; i < numNodes; i++ {
		nodes[i] = NewNode(i, network, (i+1)%numNodes)
		network.Register(i, nodes[i].Inbox)
		nodes[i].Start()
	}

	nodes[numNodes-1].SetAlive(false)

	nodes[0].startBullyElection()

	time.Sleep(10 * time.Second)

	leaderCount := 0
	var leaderID int
	for _, node := range nodes {
		if node.IsLeader {
			leaderCount++
			leaderID = node.ID
		}
	}

	if leaderCount != 1 {
		t.Errorf("Expected 1 leader, got %d", leaderCount)
	}

	if leaderID != numNodes-2 {
		t.Errorf("Expected a leader %d, got %d", numNodes-2, leaderID)
	}

	fmt.Printf("Bully election test passed. Leader is Node %d\n", leaderID)
}

func TestBullyElectionAllAlive(t *testing.T) {
	network := NewNetwork()
	numNodes := 5
	nodes := make([]*Node, numNodes)

	for i := 0; i < numNodes; i++ {
		nodes[i] = NewNode(i, network, (i+1)%numNodes)
		network.Register(i, nodes[i].Inbox)
		nodes[i].Start()
	}

	nodes[0].startBullyElection()

	time.Sleep(10 * time.Second)

	leaderCount := 0
	var leaderID int
	for _, node := range nodes {
		if node.IsLeader {
			leaderCount++
			leaderID = node.ID
		}
	}

	if leaderCount != 1 {
		t.Errorf("Expected 1 leader, got %d", leaderCount)
	}

	if leaderID != numNodes-1 {
		t.Errorf("Expected a leader %d, got %d", numNodes-1, leaderID)
	}

	fmt.Printf("Bully election all alive test passed. Leader is Node %d\n", leaderID)
}

func TestRingElection(t *testing.T) {
	network := NewNetwork()
	numNodes := 5
	nodes := make([]*Node, numNodes)

	for i := 0; i < numNodes; i++ {
		nodes[i] = NewNode(i, network, (i+1)%numNodes)
		network.Register(i, nodes[i].Inbox)
		nodes[i].Start()
	}

	// nodes[numNodes-1].SetAlive(false) // Тут так просто нельзя удалить ноду,
	// поскольку они соеденены в кольцо и если один отвалится, то колько будет нарушено

	nodes[0].StartRingElection()

	time.Sleep(10 * time.Second)

	leaderCount := 0
	var leaderID int
	for _, node := range nodes {
		if node.IsLeader {
			leaderCount++
			leaderID = node.ID
		}
	}

	if leaderCount != 1 {
		t.Errorf("Expected 1 leader, got %d", leaderCount)
	}

	if leaderID != numNodes-1 {
		t.Errorf("Expected a leader %d, got %d", numNodes-1, leaderID)
	}
	// if leaderID != numNodes-2 {
	// 	t.Errorf("Expected a leader %d, got %d", numNodes-2, leaderID)
	// }

	fmt.Printf("Ring election test passed. Leader is Node %d\n", leaderID)
}

func TestRingElectionAllAlive(t *testing.T) {
	network := NewNetwork()
	numNodes := 5
	nodes := make([]*Node, numNodes)

	for i := 0; i < numNodes; i++ {
		nodes[i] = NewNode(i, network, (i+1)%numNodes)
		network.Register(i, nodes[i].Inbox)
		nodes[i].Start()
	}

	nodes[2].StartRingElection()

	time.Sleep(10 * time.Second)

	leaderCount := 0
	var leaderID int
	for _, node := range nodes {
		if node.IsLeader {
			leaderCount++
			leaderID = node.ID
		}
	}

	if leaderCount != 1 {
		t.Errorf("Expected 1 leader, got %d", leaderCount)
	}
	if leaderID != numNodes-1 {
		t.Errorf("Expected a leader %d, got %d", numNodes-1, leaderID)
	}

	fmt.Printf("Ring election test passed. Leader is Node %d\n", leaderID)
}

func TestDataCollection(t *testing.T) {
	network := NewNetwork()
	numNodes := 5
	nodes := make([]*Node, numNodes)

	for i := 0; i < numNodes; i++ {
		nodes[i] = NewNode(i, network, (i+1)%numNodes)
		network.Register(i, nodes[i].Inbox)
		nodes[i].Start()
	}

	nodes[0].IsLeader = true
	nodes[0].LeaderID = nodes[0].ID

	nodes[0].StartGlobalCollection()

	time.Sleep(10 * time.Second)

	expectedSum := 0
	for _, node := range nodes {
		expectedSum += node.localCount
	}

	if nodes[0].collectSum != expectedSum {
		t.Errorf("Expected sum %d, got %d", expectedSum, nodes[0].collectSum)
	}

	fmt.Printf("Data collection test passed. Total sum: %d\n", nodes[0].collectSum)
}
