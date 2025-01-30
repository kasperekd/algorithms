package distributed

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type LogEntry struct {
	Term    int
	Command string
}

type RaftMessage struct {
	Kind         string
	Term         int
	From         int
	To           int
	Entries      []LogEntry
	LeaderCommit int
	Success      bool
	MatchIndex   int
	LastLogIndex int
	LastLogTerm  int
	PrevLogIndex int
	PrevLogTerm  int
}

type RaftNode struct {
	id              int
	term            int
	state           string
	log             []LogEntry
	votedFor        int
	leaderID        int
	commitIndex     int
	lastApplied     int
	peers           []*RaftNode
	inbox           chan RaftMessage
	electionTimer   *time.Timer
	heartbeatTicker *time.Ticker
	stopCh          chan struct{}
	votesReceived   int
	nextIndex       map[int]int
	matchIndex      map[int]int
}

func NewRaftNode(id int, peers []*RaftNode) *RaftNode {
	rn := &RaftNode{
		id:          id,
		term:        0,
		state:       "Follower",
		log:         []LogEntry{},
		votedFor:    -1,
		leaderID:    -1,
		commitIndex: -1,
		lastApplied: -1,
		peers:       peers,
		inbox:       make(chan RaftMessage, 100),
		stopCh:      make(chan struct{}),
		nextIndex:   make(map[int]int),
		matchIndex:  make(map[int]int),
	}
	rn.resetElectionTimer()
	return rn
}

func CreateCluster(numNodes int) []*RaftNode {
	nodes := make([]*RaftNode, numNodes)

	for i := range nodes {
		nodes[i] = NewRaftNode(i, nil)
	}

	for i := range nodes {
		peers := make([]*RaftNode, 0, numNodes-1)
		for j := range nodes {
			if i != j {
				peers = append(peers, nodes[j])
			}
		}
		nodes[i].peers = peers
	}

	return nodes
}

func StartCluster(nodes []*RaftNode) {
	for _, node := range nodes {
		go node.Run()
	}
}

func StopCluster(nodes []*RaftNode) {
	for _, node := range nodes {
		close(node.stopCh)
	}
}

func (rn *RaftNode) resetElectionTimer() {
	if rn.electionTimer != nil {
		rn.electionTimer.Stop()
	}
	timeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
	rn.electionTimer = time.AfterFunc(timeout, func() {
		if rn.state == "Follower" || rn.state == "Candidate" {
			rn.startElection()
		}
	})
}

func (rn *RaftNode) startElection() {
	rn.state = "Candidate"
	rn.term++
	rn.votedFor = rn.id
	rn.votesReceived = 1
	rn.resetElectionTimer()

	lastLogIndex := len(rn.log) - 1
	lastLogTerm := -1
	if lastLogIndex >= 0 {
		lastLogTerm = rn.log[lastLogIndex].Term
	}

	for _, peer := range rn.peers {
		if peer.id != rn.id {
			msg := RaftMessage{
				Kind:         "RequestVote",
				Term:         rn.term,
				From:         rn.id,
				To:           peer.id,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}
			peer.inbox <- msg
		}
	}
}

func (rn *RaftNode) Run() {
	for {
		select {
		case msg := <-rn.inbox:
			rn.handleMessage(msg)
		case <-rn.stopCh:
			if rn.heartbeatTicker != nil {
				rn.heartbeatTicker.Stop()
			}
			return
		}
	}
}

func (rn *RaftNode) handleMessage(msg RaftMessage) {
	switch msg.Kind {
	case "RequestVote":
		fmt.Printf("Node %d received [handleRequestVote] frome Node %d\n", rn.id, msg.From)
		rn.handleRequestVote(msg)
	case "RequestVoteResponse":
		fmt.Printf("Node %d received [handleRequestVoteResponse] frome Node %d\n", rn.id, msg.From)
		rn.handleRequestVoteResponse(msg)
	case "AppendEntries":
		fmt.Printf("Node %d received [handleAppendEntries] frome Node %d\n", rn.id, msg.From)
		rn.handleAppendEntries(msg)
	case "AppendEntriesResponse":
		fmt.Printf("Node %d received [handleAppendEntriesResponse] frome Node %d\n", rn.id, msg.From)
		rn.handleAppendEntriesResponse(msg)
	}
}

func (rn *RaftNode) handleRequestVote(msg RaftMessage) {
	if msg.Term > rn.term {
		rn.term = msg.Term
		rn.state = "Follower"
		rn.votedFor = -1
	}

	reply := RaftMessage{
		Kind: "RequestVoteResponse",
		Term: rn.term,
		From: rn.id,
		To:   msg.From,
	}

	if msg.Term < rn.term {
		reply.Success = false
	} else if (rn.votedFor == -1 || rn.votedFor == msg.From) && rn.isLogUpToDate(msg.LastLogTerm, msg.LastLogIndex) {
		rn.votedFor = msg.From
		reply.Success = true
		rn.resetElectionTimer()
	} else {
		reply.Success = false
	}

	rn.sendMessage(reply)
}

func (rn *RaftNode) isLogUpToDate(candidateTerm, candidateIndex int) bool {
	lastIndex := len(rn.log) - 1
	if lastIndex == -1 {
		return true
	}
	lastTerm := rn.log[lastIndex].Term
	return candidateTerm > lastTerm || (candidateTerm == lastTerm && candidateIndex >= lastIndex)
}

func (rn *RaftNode) handleRequestVoteResponse(msg RaftMessage) {
	if rn.state != "Candidate" || msg.Term != rn.term {
		return
	}

	if msg.Success {
		rn.votesReceived++
		if rn.votesReceived > len(rn.peers)/2 {
			rn.becomeLeader()
		}
	}
}

func (rn *RaftNode) becomeLeader() {
	fmt.Printf("Node %d became a leader !\n", rn.id)

	rn.state = "Leader"
	rn.leaderID = rn.id
	for _, peer := range rn.peers {
		rn.nextIndex[peer.id] = len(rn.log)
		rn.matchIndex[peer.id] = -1
	}
	rn.heartbeatTicker = time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-rn.heartbeatTicker.C:
				rn.sendHeartbeats()
			case <-rn.stopCh:
				return
			}
		}
	}()
}

func (rn *RaftNode) sendHeartbeats() {
	for _, peer := range rn.peers {
		if peer.id == rn.id {
			continue
		}

		prevLogIndex := rn.nextIndex[peer.id] - 1
		prevLogTerm := -1
		if prevLogIndex >= 0 {
			prevLogTerm = rn.log[prevLogIndex].Term
		}

		entries := rn.log[rn.nextIndex[peer.id]:]
		msg := RaftMessage{
			Kind:         "AppendEntries",
			Term:         rn.term,
			From:         rn.id,
			To:           peer.id,
			Entries:      entries,
			LeaderCommit: rn.commitIndex,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
		}
		peer.inbox <- msg
	}
}

func (rn *RaftNode) handleAppendEntries(msg RaftMessage) {
	reply := RaftMessage{
		Kind: "AppendEntriesResponse",
		Term: rn.term,
		From: rn.id,
		To:   msg.From,
	}

	if msg.Term < rn.term {
		reply.Success = false
		rn.sendMessage(reply)
		return
	}

	rn.resetElectionTimer()
	if msg.Term > rn.term {
		rn.term = msg.Term
		rn.state = "Follower"
		rn.votedFor = -1
	}

	if msg.PrevLogIndex != -1 && (len(rn.log) <= msg.PrevLogIndex || rn.log[msg.PrevLogIndex].Term != msg.PrevLogTerm) {
		reply.Success = false
		reply.MatchIndex = len(rn.log) - 1
	} else {
		rn.log = append(rn.log[:msg.PrevLogIndex+1], msg.Entries...)
		reply.Success = true
		reply.MatchIndex = len(rn.log) - 1

		if msg.LeaderCommit > rn.commitIndex {
			rn.commitIndex = min(msg.LeaderCommit, len(rn.log)-1)
			rn.applyLogs()
		}
	}

	rn.sendMessage(reply)
}

func (rn *RaftNode) applyLogs() {
	for rn.lastApplied < rn.commitIndex {
		rn.lastApplied++
		entry := rn.log[rn.lastApplied]
		fmt.Printf("Node %d applied command '%s' at index %d\n", rn.id, entry.Command, rn.lastApplied)
	}
}

func (rn *RaftNode) handleAppendEntriesResponse(msg RaftMessage) {
	if rn.state != "Leader" || msg.Term != rn.term {
		return
	}

	if msg.Success {
		rn.matchIndex[msg.From] = msg.MatchIndex
		rn.nextIndex[msg.From] = msg.MatchIndex + 1

		matchIndices := make([]int, 0, len(rn.peers))
		for _, idx := range rn.matchIndex {
			matchIndices = append(matchIndices, idx)
		}
		matchIndices = append(matchIndices, len(rn.log)-1)
		sort.Ints(matchIndices)
		N := matchIndices[len(matchIndices)/2]

		if N > rn.commitIndex && (N < len(rn.log) && rn.log[N].Term == rn.term) {
			rn.commitIndex = N
			rn.applyLogs()
		}
	} else {
		rn.nextIndex[msg.From] = max(1, rn.nextIndex[msg.From]-1)
		prevLogIndex := rn.nextIndex[msg.From] - 1
		prevLogTerm := -1
		if prevLogIndex >= 0 {
			prevLogTerm = rn.log[prevLogIndex].Term
		}
		entries := rn.log[rn.nextIndex[msg.From]:]

		msg := RaftMessage{
			Kind:         "AppendEntries",
			Term:         rn.term,
			From:         rn.id,
			To:           msg.From,
			Entries:      entries,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			LeaderCommit: rn.commitIndex,
		}
		rn.sendMessage(msg)
	}
}

func (rn *RaftNode) sendMessage(msg RaftMessage) {
	for _, peer := range rn.peers {
		if peer.id == msg.To {
			peer.inbox <- msg
			return
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func findLeader(nodes []*RaftNode) *RaftNode {
	for _, node := range nodes {
		if node.state == "Leader" {
			return node
		}
	}
	return nil
}

func Task8() {
	cluster := CreateCluster(3)
	StartCluster(cluster)

	time.Sleep(500 * time.Millisecond)

	if leader := findLeader(cluster); leader != nil {
		leader.log = append(leader.log, LogEntry{Term: leader.term, Command: "X"})
		fmt.Printf("Leader %d added command 'X'\n", leader.id)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("\nSimulating leader failure")
	if leader := findLeader(cluster); leader != nil {
		close(leader.stopCh)
		leader.state = "Follower"
	}

	time.Sleep(1 * time.Second)

	if newLeader := findLeader(cluster); newLeader != nil {
		newLeader.log = append(newLeader.log, LogEntry{Term: newLeader.term, Command: "Y"})
		fmt.Printf("New leader %d added command 'Y'\n", newLeader.id)
	}
	time.Sleep(2 * time.Second)

	// StopCluster(cluster)

}
