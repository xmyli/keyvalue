package raft

import (
	"net/rpc"
)

func (rf *Raft) toLeader() {
	rf.mu.Lock()

	rf.role = Leader

	rf.nextIndex = []int{}
	rf.matchIndex = []int{}
	for range rf.peers {
		rf.nextIndex = append(rf.nextIndex, len(rf.log)-1)
		rf.matchIndex = append(rf.matchIndex, 0)
	}

	rf.mu.Unlock()

	go rf.updateFollowers()
	go rf.updateCommitIndex()
}

func (rf *Raft) updateFollowers() {
	for {
		rf.mu.Lock()

		if rf.role != Leader {
			rf.mu.Unlock()
			return
		}

		for index, peer := range rf.peers {
			if index == rf.me {
				continue
			}
			go rf.appendEntries(index, peer)
		}

		rf.mu.Unlock()

		sleepFor(20, 20)
	}
}

func (rf *Raft) appendEntries(index int, peer string) {
	for {
		rf.mu.Lock()

		if rf.role != Leader {
			rf.mu.Unlock()
			break
		}

		nextIndex := rf.nextIndex[index]

		args := AppendEntriesArgs{}
		reply := AppendEntriesReply{}

		args.Term = rf.currentTerm
		args.LeaderId = rf.me
		if len(rf.log) >= nextIndex && nextIndex > 0 {
			args.PrevLogIndex = rf.log[nextIndex-1].Index
			args.PrevLogTerm = rf.log[nextIndex-1].Term
		} else {
			args.PrevLogIndex = -1
			args.PrevLogTerm = -1
		}
		args.Entries = rf.log[nextIndex:]
		args.LeaderCommit = rf.commitIndex

		rf.mu.Unlock()

		rpcClient, err := rpc.Dial("tcp", peer)
		if err != nil {
			break
		}

		err = rpcClient.Call("Raft.AppendEntries", &args, &reply)
		if err != nil {
			break
		}

		rpcClient.Close()

		rf.mu.Lock()

		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			if rf.role == Leader {
				go rf.toFollower()
			}
			rf.mu.Unlock()
			break
		}

		if len(args.Entries) == 0 {
			rf.mu.Unlock()
			break
		}

		if rf.role != Leader {
			rf.mu.Unlock()
			break
		}

		if reply.Success {
			rf.nextIndex[index] = args.Entries[len(args.Entries)-1].Index + 1
			rf.matchIndex[index] = args.Entries[len(args.Entries)-1].Index
			rf.mu.Unlock()
			break
		}

		rf.nextIndex[index]--

		rf.mu.Unlock()
	}
}

func (rf *Raft) updateCommitIndex() {
	for {
		rf.mu.Lock()

		if rf.role != Leader {
			rf.mu.Unlock()
			return
		}

		for N := rf.commitIndex + 1; isMajorityGreaterOrEqual(N, rf.matchIndex); N++ {
			if rf.log[N].Term == rf.currentTerm {
				rf.commitIndex = N
			}
		}

		rf.mu.Unlock()

		sleepFor(20, 20)
	}
}

func (rf *Raft) SendCommand(command interface{}) (int, int, bool) {
	rf.mu.Lock()

	if rf.role != Leader {
		rf.mu.Unlock()
		return -1, -1, false
	}

	entry := Entry{}
	entry.Command = command
	entry.Index = len(rf.log)
	entry.Term = rf.currentTerm

	rf.log = append(rf.log, entry)

	rf.mu.Unlock()

	return entry.Index, entry.Term, true
}
