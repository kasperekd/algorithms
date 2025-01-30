package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/kasperekd/algorithms/distributed"
// )

// func main() {
// 	network := distributed.NewNetwork()
// 	nodes := make([]*distributed.Node, 0)
// 	numNodes := 5

// 	// Создаем узлы и настраиваем кольцо
// 	for i := 1; i <= numNodes; i++ {
// 		nextID := i%numNodes + 1
// 		node := distributed.NewNode(i, network, nextID)
// 		node.LeaderID = -1 // Изначально лидер неизвестен
// 		node.SetAlive(true)
// 		network.Register(i, node.Inbox)
// 		nodes = append(nodes, node)
// 		node.Start()
// 	}

// 	// Даем время на инициализацию
// 	time.Sleep(time.Second)

// 	// Узел 1 запускает выборы
// 	fmt.Println("--- Запуск выборов ---")
// 	nodes[0].StartBullyElection()
// 	time.Sleep(2 * time.Second)

// 	// Определяем лидера
// 	var leader *distributed.Node
// 	for _, n := range nodes {
// 		if n.IsLeader {
// 			leader = n
// 			break
// 		}
// 	}
// 	fmt.Printf("Лидер выбран: Узел %d\n\n", leader.ID)

// 	// Сбор данных лидером
// 	fmt.Println("--- Сбор данных ---")
// 	leader.StartGlobalCollection()
// 	time.Sleep(3 * time.Second)

// 	// Имитация сбоя узла 2
// 	fmt.Println("\n--- Имитация сбоя Узла 2 ---")
// 	nodes[1].SetAlive(false)
// 	fmt.Println("Узел 2 недоступен")

// 	// Повторный сбор данных
// 	fmt.Println("\n--- Повторный сбор данных ---")
// 	leader.StartGlobalCollection()
// 	time.Sleep(11 * time.Second) // Ожидаем таймаут

// 	// Имитация сбоя лидера и новые выборы
// 	fmt.Println("\n--- Сбой лидера и новые выборы ---")
// 	leader.SetAlive(false)
// 	fmt.Printf("Лидер (Узел %d) недоступен\n", leader.ID)
// 	nodes[0].StartBullyElection()
// 	time.Sleep(20 * time.Second)

// 	// Поиск нового лидера
// 	var newLeader *distributed.Node
// 	for _, n := range nodes {
// 		if n.IsLeader && n.Alive {
// 			newLeader = n
// 			break
// 		}
// 	}
// 	if newLeader != nil {
// 		fmt.Printf("Новый лидер: Узел %d\n\n", newLeader.ID)
// 	} else {
// 		fmt.Println("Лидер не выбран")
// 		return
// 	}

// 	// Сбор данных новым лидером
// 	fmt.Println("--- Сбор данных новым лидером ---")
// 	newLeader.StartGlobalCollection()
// 	time.Sleep(5 * time.Second)
// }
