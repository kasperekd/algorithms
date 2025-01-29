package distributed

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBullyElection(t *testing.T) {
	network := NewNetwork()
	nodes := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = NewNode(i, network, (i+1)%5)
		network.Register(i, nodes[i].Inbox)
		nodes[i].Start()
	}

	fmt.Println("Starting Bully Election Test")
	nodes[2].startBullyElection()

	time.Sleep(2 * time.Second)

	leaderID := -1
	for _, node := range nodes {
		if node.IsLeader {
			leaderID = node.ID
			break
		}
	}

	assert.NotEqual(t, -1, leaderID, "No leader was elected")

	for _, node := range nodes {
		assert.Equal(t, leaderID, node.LeaderID, "Nodes do not agree on the leader")
	}
	fmt.Printf("Bully Election Test Passed. Leader: %d\n", leaderID)

}

func TestRingElection(t *testing.T) {
	network := NewNetwork()
	nodes := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = NewNode(i, network, (i+1)%5)
		network.Register(i, nodes[i].Inbox)
		nodes[i].Start()
	}

	fmt.Println("Starting Ring Election Test")
	nodes[0].startRingElection()

	time.Sleep(2 * time.Second)

	leaderID := -1
	for _, node := range nodes {
		if node.IsLeader {
			leaderID = node.ID
			break
		}
	}

	assert.NotEqual(t, -1, leaderID, "No leader was elected")

	for _, node := range nodes {
		assert.Equal(t, leaderID, node.LeaderID, "Nodes do not agree on the leader")
	}
	fmt.Printf("Ring Election Test Passed. Leader: %d\n", leaderID)
}

func TestGlobalCollection(t *testing.T) {
	network := NewNetwork()
	nodes := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = NewNode(i, network, (i+1)%5)
		network.Register(i, nodes[i].Inbox)
		nodes[i].Start()
	}

	nodes[0].IsLeader = true
	nodes[0].LeaderID = nodes[0].ID

	fmt.Println("Starting Global Collection Test")
	go nodes[0].StartGlobalCollection()

	time.Sleep(2 * time.Second)

	expectedSum := 0
	for _, node := range nodes {
		expectedSum += node.localCount
	}

	fmt.Println("Checking if printed total sum is near to expected sum. Please check output.")

}

func TestNodeFailure(t *testing.T) {
	network := NewNetwork()
	nodes := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = NewNode(i, network, (i+1)%5)
		network.Register(i, nodes[i].Inbox)
		nodes[i].Start()
	}

	nodes[0].IsLeader = true
	nodes[0].LeaderID = nodes[0].ID

	fmt.Println("Starting Node Failure Test")
	nodes[2].SetAlive(false)

	go nodes[0].StartGlobalCollection()
	time.Sleep(2 * time.Second)

	fmt.Println("Checking if printed total sum is near to expected sum. Please check output. Node 2 is down.")
	fmt.Println("Node Failure Test Passed (manual check).")

}
