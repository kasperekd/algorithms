package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/kasperekd/algorithms/distributed"
)

func main() {
	network := distributed.NewNetwork()
	nodeIDs := []int{0, 1, 2, 3}
	nodes := make(map[int]*distributed.Node)

	for i, id := range nodeIDs {
		nextID := nodeIDs[(i+1)%len(nodeIDs)]
		node := distributed.NewNode(id, network, nextID)
		node.LeaderID = -1
		nodes[id] = node
		network.Register(id, node.Inbox)
	}

	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(n *distributed.Node) {
			defer wg.Done()
			n.Start()
		}(node)
	}

	fmt.Println("=== Инициализация узлов завершена ===")

	fmt.Println("\n=== Узел 1 запускает кольцевые выборы ===")
	nodes[1].StartRingElection()
	time.Sleep(2 * time.Second)

	var leader *distributed.Node
	for _, node := range nodes {
		if node.IsLeader {
			leader = node
			break
		}
	}
	if leader != nil {
		fmt.Printf("\nЛидер %d выбран\n", leader.ID)
		fmt.Println("Запуск сбора данных...")
		leader.StartGlobalCollection()
		time.Sleep(3 * time.Second)
	} else {
		fmt.Println("Лидер не найден")
	}

	fmt.Println("\n=== Имитация сбоя узла 3 ===")
	nodes[3].SetAlive(false)

	fmt.Println("\n=== Перезапуск выборов ===")
	nodes[0].StartRingElection()
	time.Sleep(2 * time.Second)

	var newLeader *distributed.Node
	for id, node := range nodes {
		if id == 3 {
			continue
		}
		if node.IsLeader {
			newLeader = node
			break
		}
	}
	if newLeader != nil {
		fmt.Printf("\nНовый лидер %d выбран\n", newLeader.ID)
		fmt.Println("Повторный сбор данных...")
		newLeader.StartGlobalCollection()
		time.Sleep(3 * time.Second)
	} else {
		fmt.Println("Новый лидер не выбран")
	}

	wg.Wait()
}
