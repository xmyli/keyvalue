package raft

import (
	"time"
)

type AppendEntriesArgs struct {
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []Entry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	rf.mu.Lock()

	rf.lastHeartbeatTime = time.Now().UnixMilli()

	reply.Term = rf.currentTerm

	if args.Term < rf.currentTerm {
		reply.Success = false
		rf.mu.Unlock()
		return nil
	} else if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term

		if rf.role != Follower {
			rf.mu.Unlock()
			rf.toFollower()
			rf.mu.Lock()
		}
	} else {
		if rf.role == Candidate {
			rf.mu.Unlock()
			rf.toFollower()
			rf.mu.Lock()
		}
	}

	if args.PrevLogIndex < len(rf.log) {
		if args.PrevLogIndex != -1 && args.PrevLogTerm != rf.log[args.PrevLogIndex].Term {
			reply.Success = false
			rf.mu.Unlock()
			return nil
		}
	} else {
		reply.Success = false
		rf.mu.Unlock()
		return nil
	}

	for _, entry := range args.Entries {
		if entry.Index < len(rf.log) && entry.Term != rf.log[entry.Index].Term {
			new := truncateLog(rf.log, entry.Index)
			rf.log = new
			break
		}
	}

	for _, entry := range args.Entries {
		if entry.Index >= len(rf.log) {
			rf.log = append(rf.log, entry)
		}
	}

	if args.LeaderCommit > rf.commitIndex {
		if args.LeaderCommit <= len(rf.log)-1 {
			rf.commitIndex = args.LeaderCommit
		} else {
			rf.commitIndex = len(rf.log) - 1
		}
	}

	reply.Success = true

	rf.mu.Unlock()

	return nil
}

func truncateLog(log []Entry, to int) []Entry {
	return log[:to]
}
