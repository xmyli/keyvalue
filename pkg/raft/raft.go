package raft

import (
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Role string

const (
	Follower  Role = "FOLLOWER"
	Candidate Role = "CANDIDATE"
	Leader    Role = "LEADER"
)

type Entry struct {
	Command interface{}
	Index   int
	Term    int
}

type Raft struct {
	port string

	mu    sync.Mutex
	peers []string
	me    int

	commandsCh chan interface{}

	currentTerm int
	votedFor    int
	log         []Entry

	role        Role
	commitIndex int
	lastApplied int

	nextIndex  []int
	matchIndex []int

	lastHeartbeatTime int64
	votes             int
}

func (rf *Raft) IsLeader() bool {
	rf.mu.Lock()
	isLeader := rf.role == Leader
	rf.mu.Unlock()

	return isLeader
}

func (rf *Raft) Start() {
	rpcs := rpc.NewServer()
	rpcs.Register(rf)
	listener, err := net.Listen("tcp", ":"+rf.port)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				go rpcs.ServeConn(conn)
			} else {
				break
			}
		}
		listener.Close()
	}()
}

func CreateRaftServer(port string, peers []string, me int, commandsCh chan interface{}) *Raft {
	rf := &Raft{}

	rf.port = port

	rf.peers = peers
	rf.me = me

	rf.commandsCh = commandsCh

	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.log = []Entry{}

	rf.log = append(rf.log, Entry{nil, 0, 0})

	rf.nextIndex = []int{}
	rf.matchIndex = []int{}

	for range peers {
		rf.nextIndex = append(rf.nextIndex, 0)
		rf.matchIndex = append(rf.matchIndex, 0)
	}

	rf.votedFor = -1

	go rf.toFollower()
	go rf.appendCommands()

	return rf
}

func (rf *Raft) appendCommands() {
	for {
		rf.mu.Lock()
		if rf.commitIndex > rf.lastApplied {
			rf.lastApplied++

			command := rf.log[rf.lastApplied].Command

			rf.mu.Unlock()

			rf.commandsCh <- command
		} else {
			rf.mu.Unlock()
		}
		sleepFor(20, 20)
	}
}
