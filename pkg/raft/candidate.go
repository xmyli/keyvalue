package raft

import (
	"net/rpc"
)

func (rf *Raft) toCandidate() {
	rf.mu.Lock()
	rf.role = Candidate
	rf.mu.Unlock()

	go rf.doElection()
}

func (rf *Raft) doElection() {
	for {
		rf.mu.Lock()

		if rf.role != Candidate {
			rf.mu.Unlock()
			return
		}

		rf.currentTerm++
		rf.votedFor = rf.me
		rf.votes = 1

		for index, peer := range rf.peers {
			if index == rf.me {
				continue
			}
			go rf.requestVote(index, peer)
		}

		rf.mu.Unlock()

		sleepFor(150, 300)
	}
}

func (rf *Raft) requestVote(index int, peer string) {
	rf.mu.Lock()

	if rf.role != Candidate {
		rf.mu.Unlock()
		return
	}

	args := RequestVoteArgs{}
	reply := RequestVoteReply{}

	args.Term = rf.currentTerm
	args.CandidateId = rf.me
	args.LastLogIndex = len(rf.log) - 1
	args.LastLogTerm = rf.log[len(rf.log)-1].Term

	rf.mu.Unlock()

	rpcClient, err := rpc.Dial("tcp", peer)
	if err != nil {
		return
	}

	err = rpcClient.Call("Raft.RequestVote", &args, &reply)
	if err != nil {
		return
	}

	rpcClient.Close()

	rf.mu.Lock()

	if rf.role != Candidate {
		rf.mu.Unlock()
		return
	}

	if reply.Term > rf.currentTerm {
		rf.currentTerm = reply.Term

		if rf.role == Candidate {
			go rf.toFollower()
		}
	}

	if reply.Term == rf.currentTerm && reply.VoteGranted {
		rf.votes++

		if rf.votes > len(rf.peers)/2 {
			go rf.toLeader()
		}
	}

	rf.mu.Unlock()
}
