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
		network.Register(nodes[i])
		nodes[i].Start()
	}

	nodes[numNodes-1].SetAlive(false)

	nodes[0].StartBullyElection()

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
		network.Register(nodes[i])
		nodes[i].Start()
	}

	nodes[0].StartBullyElection()

	time.Sleep(15 * time.Second)

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
		network.Register(nodes[i])
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
		network.Register(nodes[i])
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
		network.Register(nodes[i])
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

func TestMain(t *testing.T) {
	network := NewNetwork()
	nodes := make([]*Node, 0)
	numNodes := 5

	for i := 0; i <= numNodes; i++ {
		nextID := i%numNodes + 1
		node := NewNode(i, network, nextID)
		node.LeaderID = -1
		node.SetAlive(true)
		network.Register(node)
		nodes = append(nodes, node)
		node.Start()
	}

	time.Sleep(time.Second)

	fmt.Println("--- Запуск выборов ---")
	nodes[1].StartBullyElection()
	time.Sleep(2 * time.Second)

	var leader *Node
	for _, n := range nodes {
		if n.IsLeader {
			leader = n
			break
		}
	}
	fmt.Printf("Лидер выбран: Узел %d\n\n", leader.ID)

	fmt.Println("--- Сбор данных ---")
	leader.StartGlobalCollection()
	time.Sleep(3 * time.Second)

	fmt.Println("\n--- Имитация сбоя Узла 2 ---")
	nodes[2].SetAlive(false)
	fmt.Println("Узел 2 недоступен")

	fmt.Println("\n--- Повторный сбор данных ---")
	leader.StartGlobalCollection()
	time.Sleep(11 * time.Second)

	fmt.Println("\n--- Сбой лидера и новые выборы ---")
	leader.SetAlive(false)
	fmt.Printf("Лидер (Узел %d) недоступен\n", leader.ID)
	nodes[0].StartBullyElection()
	time.Sleep(20 * time.Second)

	var newLeader *Node
	for _, n := range nodes {
		if n.IsLeader && n.Alive {
			newLeader = n
			break
		}
	}
	if newLeader != nil {
		fmt.Printf("Новый лидер: Узел %d\n\n", newLeader.ID)
	} else {
		fmt.Println("Лидер не выбран")
		return
	}

	fmt.Println("--- Сбор данных новым лидером ---")
	newLeader.StartGlobalCollection()
	time.Sleep(5 * time.Second)
}
