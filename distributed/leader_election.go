package distributed

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

type Message struct {
	Kind      string
	Algorithm string // bully или ring
	FromID    int
	ToID      int
	Data      interface{}
	MaxID     int
	OriginID  int
}

type Network struct {
	mu       sync.RWMutex
	channels map[int]chan Message
	nodes    map[int]*Node
}

type Node struct {
	ID            int
	IsLeader      bool
	Alive         bool
	LeaderID      int
	NextID        int
	Inbox         chan Message
	network       *Network
	localCount    int
	mu            sync.Mutex
	electionTimer *time.Timer
	collectSum    int
	collectCount  int
	collectDone   chan struct{}
}

func NewNode(id int, network *Network, nextID int) *Node {
	rand.Seed(uint64(time.Now().UnixNano()))
	return &Node{
		ID:          id,
		Alive:       true,
		Inbox:       make(chan Message, 100),
		network:     network,
		NextID:      nextID,
		localCount:  rand.Intn(50) + 50,
		collectDone: make(chan struct{}),
	}
}

func (n *Node) Start() {
	go n.processMessages()
}

func (n *Node) processMessages() {
	for msg := range n.Inbox {
		if !n.Alive {
			continue
		}
		switch msg.Kind {
		case "ELECTION":
			if msg.Algorithm == "ring" {
				n.handleRingElection(msg)
			} else {
				n.handleElection(msg)
			}
		case "OK":
			n.handleOk(msg)
		case "COORDINATOR":
			n.handleCoordinator(msg)
		case "COLLECT":
			n.handleCollect(msg)
		case "COLLECT_REPLY":
			n.handleCollectReply(msg)
		}
	}
}

func (n *Node) SetAlive(alive bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Alive = alive
}

func (n *Network) IsAlive(id int) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	node, exists := n.nodes[id]
	return exists && node.Alive
}

func (n *Network) GetAliveNodes() []int {
	n.mu.RLock()
	defer n.mu.RUnlock()

	alive := make([]int, 0)
	for id, node := range n.nodes {
		node.mu.Lock()
		if node.Alive {
			alive = append(alive, id)
		}
		node.mu.Unlock()
	}
	return alive
}

func NewNetwork() *Network {
	return &Network{
		channels: make(map[int]chan Message),
		nodes:    make(map[int]*Node),
	}
}

func (n *Network) Register(node *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.channels[node.ID] = node.Inbox
	n.nodes[node.ID] = node
}
func (n *Network) Send(toID int, msg Message) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	ch, ok := n.channels[toID]
	if !ok {
		return false
	}
	select {
	case ch <- msg:
		return true
	default:
		return false
	}
}

func (n *Network) GetNodes() []int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	ids := make([]int, 0, len(n.channels))
	for id := range n.channels {
		ids = append(ids, id)
	}
	return ids
}

func (n *Node) StartBullyElection() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.electionTimer != nil {
		return
	}

	fmt.Printf("Node %d starting Bully election\n", n.ID)

	higherNodes := make([]int, 0)
	for _, id := range n.network.GetNodes() {
		if id > n.ID {
			higherNodes = append(higherNodes, id)
		}
	}

	if len(higherNodes) == 0 {
		n.declareLeader()
		return
	}

	for _, id := range higherNodes {
		fmt.Printf("Node %d sending ELECTION to Node %d\n", n.ID, id)
		n.network.Send(id, Message{
			Kind:   "ELECTION",
			FromID: n.ID,
			ToID:   id,
		})
	}

	n.electionTimer = time.AfterFunc(10*time.Second, func() {
		n.mu.Lock()
		defer n.mu.Unlock()

		if n.electionTimer != nil {
			n.declareLeader()
			n.electionTimer = nil
		}
	})
}

func (n *Node) handleElection(msg Message) {
	fmt.Printf("Node %d received ELECTION message from Node %d\n", n.ID, msg.FromID)

	n.mu.Lock()
	defer n.mu.Unlock()

	if msg.FromID < n.ID {
		go func() {
			n.network.Send(msg.FromID, Message{Kind: "OK", FromID: n.ID, ToID: msg.FromID})
		}()

		if n.electionTimer == nil {
			go n.StartBullyElection()
		}
	} else if msg.FromID > n.ID {
		if n.electionTimer != nil {
			n.electionTimer.Stop()
			n.electionTimer = nil
		}
	}
}

func (n *Node) handleOk(msg Message) {
	fmt.Printf("Node %d received OK message from Node %d\n", n.ID, msg.FromID)
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.electionTimer != nil {
		n.electionTimer.Stop()
		n.electionTimer = nil
	}
}

func (n *Node) declareLeader() {
	n.IsLeader = true
	n.LeaderID = n.ID
	fmt.Printf("Node %d declared itself as leader\n", n.ID)
	allNodes := n.network.GetNodes()
	for _, id := range allNodes {
		if id != n.ID {
			n.network.Send(id, Message{Kind: "COORDINATOR", FromID: n.ID, ToID: id})
		}
	}
}

func (n *Node) handleCoordinator(msg Message) {
	if msg.Algorithm == "ring" {
		n.mu.Lock()
		defer n.mu.Unlock()

		fmt.Printf("Node %d received RING coordinator: Leader %d\n",
			n.ID, msg.MaxID)

		n.LeaderID = msg.MaxID
		n.IsLeader = (n.ID == msg.MaxID)

		if msg.OriginID != n.ID {
			forwardMsg := msg
			forwardMsg.ToID = n.NextID
			forwardMsg.FromID = n.ID
			n.network.Send(n.NextID, forwardMsg)
		}
		return
	}

	n.LeaderID = msg.FromID
	n.IsLeader = (n.ID == msg.FromID)
	fmt.Printf("Node %d received COORDINATOR message from Node %d. Leader is now %d\n", n.ID, msg.FromID, n.LeaderID)
}

func (n *Node) StartRingElection() {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Printf("Node %d initiating RING election\n", n.ID)
	msg := Message{
		Kind:      "ELECTION",
		Algorithm: "ring",
		FromID:    n.ID,
		ToID:      n.NextID,
		OriginID:  n.ID,
		MaxID:     n.ID,
	}
	n.network.Send(n.NextID, msg)
}

func (n *Node) handleRingElection(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()

	currentMax := msg.MaxID
	if n.Alive && n.ID > currentMax {
		currentMax = n.ID
	}

	if msg.OriginID == n.ID {
		if nextID := n.findNextAlive(); nextID != -1 {
			n.network.Send(nextID, Message{
				Kind:      "COORDINATOR",
				Algorithm: "ring",
				MaxID:     currentMax,
				OriginID:  n.ID,
			})
		} else {
			n.declareLeader()
		}
		return
	}

	if nextID := n.findNextAlive(); nextID != -1 {
		n.network.Send(nextID, Message{
			Kind:      "ELECTION",
			Algorithm: "ring",
			OriginID:  msg.OriginID,
			MaxID:     currentMax,
		})
	}
}

func (n *Node) StartGlobalCollection() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.IsLeader || !n.Alive {
		return
	}

	aliveNodes := n.network.GetAliveNodes()
	expected := len(aliveNodes) - 1 // Исключаем себя
	n.collectSum = 0
	n.collectCount = 0
	n.collectDone = make(chan struct{})

	for _, id := range aliveNodes {
		if id != n.ID {
			n.network.Send(id, Message{Kind: "COLLECT", FromID: n.ID, ToID: id})
		}
	}

	go func() {
		select {
		case <-n.collectDone:
			n.mu.Lock()
			total := n.collectSum + n.localCount
			fmt.Printf("Total sum from %d/%d nodes: %d\n",
				n.collectCount+1, expected+1, total)
			n.mu.Unlock()
		case <-time.After(5 * time.Second):
			n.mu.Lock()
			fmt.Printf("Partial data: %d/%d (Alive: %v)\n",
				n.collectCount, expected, aliveNodes)
			n.mu.Unlock()
		}
	}()
}

func (n *Node) handleCollectReply(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if data, ok := msg.Data.(int); ok {
		n.collectSum += data
		n.collectCount++
		if n.collectCount >= len(n.network.GetAliveNodes())-1 {
			close(n.collectDone)
		}
	}
}

func (n *Node) findNextAlive() int {
	all := n.network.GetNodes()
	for i := 0; i < len(all); i++ {
		candidate := (n.NextID + i) % len(all)
		if n.network.IsAlive(candidate) && candidate != n.ID {
			return candidate
		}
	}
	return -1
}

func (n *Node) handleCollect(msg Message) {
	reply := Message{
		Kind:   "COLLECT_REPLY",
		FromID: n.ID,
		ToID:   msg.FromID,
		Data:   n.localCount,
	}
	n.network.Send(msg.FromID, reply)
}

// func (n *Node) handleCollectReply(msg Message) {
// 	n.mu.Lock()
// 	defer n.mu.Unlock()
// 	if data, ok := msg.Data.(int); ok {
// 		n.collectSum += data
// 		n.collectCount++
// 		if n.collectCount == len(n.network.GetNodes())-1 {
// 			close(n.collectDone)
// 		}
// 	}
// }
