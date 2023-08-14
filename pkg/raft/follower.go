package raft

import (
	"time"
)

func (rf *Raft) toFollower() {
	rf.mu.Lock()
	rf.role = Follower
	rf.votedFor = -1
	rf.mu.Unlock()

	go rf.checkHeartbeat()
}

func (rf *Raft) checkHeartbeat() {
	for {
		sleepFor(150, 300)

		rf.mu.Lock()

		if rf.role != Follower {
			rf.mu.Unlock()
			return
		}

		if rf.role != Follower {
			rf.mu.Unlock()
			return
		}

		diff := time.Now().UnixMilli() - rf.lastHeartbeatTime
		if diff > 150 {
			if rf.role == Follower {
				go rf.toCandidate()
			}
		}

		rf.mu.Unlock()
	}
}
