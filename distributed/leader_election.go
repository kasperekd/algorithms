package distributed

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

type Message struct {
	Kind     string
	FromID   int
	ToID     int
	Data     interface{}
	MaxID    int
	OriginID int
}

type Network struct {
	mu       sync.RWMutex
	channels map[int]chan Message
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
			n.handleElection(msg)
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

func NewNetwork() *Network {
	return &Network{
		channels: make(map[int]chan Message),
	}
}

func (n *Network) Register(id int, ch chan Message) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.channels[id] = ch
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

func (n *Node) startBullyElection() {
	n.mu.Lock()
	defer n.mu.Unlock()

	allNodes := n.network.GetNodes()
	var higherNodes []int
	for _, id := range allNodes {
		if id > n.ID {
			higherNodes = append(higherNodes, id)
		}
	}

	if len(higherNodes) == 0 {
		n.declareLeader()
		return
	}

	for _, id := range higherNodes {
		n.network.Send(id, Message{Kind: "ELECTION", FromID: n.ID, ToID: id})
	}

	n.electionTimer = time.AfterFunc(5*time.Second, func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		n.declareLeader()
	})
}

func (n *Node) handleElection(msg Message) {
	if msg.FromID < n.ID {
		n.network.Send(msg.FromID, Message{Kind: "OK", FromID: n.ID, ToID: msg.FromID})
		n.startBullyElection()
	}
}

func (n *Node) handleOk(msg Message) {
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
	allNodes := n.network.GetNodes()
	for _, id := range allNodes {
		if id != n.ID {
			n.network.Send(id, Message{Kind: "COORDINATOR", FromID: n.ID, ToID: id})
		}
	}
}

func (n *Node) handleCoordinator(msg Message) {
	n.LeaderID = msg.FromID
	n.IsLeader = (n.ID == msg.FromID)
}

func (n *Node) startRingElection() {
	msg := Message{
		Kind:     "ELECTION",
		FromID:   n.ID,
		OriginID: n.ID,
		MaxID:    n.ID,
	}
	n.network.Send(n.NextID, msg)
}

func (n *Node) handleRingElection(msg Message) {
	if msg.MaxID < n.ID {
		msg.MaxID = n.ID
	}
	msg.FromID = n.ID
	n.network.Send(n.NextID, msg)

	if msg.OriginID == n.ID {
		n.network.Send(n.NextID, Message{
			Kind:   "COORDINATOR",
			FromID: n.ID,
			MaxID:  msg.MaxID,
		})
	}
}

func (n *Node) StartGlobalCollection() {
	if !n.IsLeader {
		return
	}

	n.mu.Lock()
	n.collectSum = 0
	n.collectCount = 0
	n.collectDone = make(chan struct{})
	allNodes := n.network.GetNodes()
	// expected := len(allNodes) - 1
	n.mu.Unlock()

	for _, id := range allNodes {
		if id != n.ID {
			n.network.Send(id, Message{Kind: "COLLECT", FromID: n.ID, ToID: id})
		}
	}

	select {
	case <-n.collectDone:
		n.mu.Lock()
		fmt.Printf("Total sum: %d\n", n.collectSum)
		n.mu.Unlock()
	case <-time.After(10 * time.Second):
		fmt.Println("Timeout during data collection")
	}
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

func (n *Node) handleCollectReply(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if data, ok := msg.Data.(int); ok {
		n.collectSum += data
		n.collectCount++
		if n.collectCount == len(n.network.GetNodes())-1 {
			close(n.collectDone)
		}
	}
}
