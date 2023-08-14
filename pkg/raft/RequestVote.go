package raft

type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	rf.mu.Lock()

	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		rf.mu.Unlock()
		return nil
	}

	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term

		if rf.role != Follower {
			rf.mu.Unlock()
			rf.toFollower()
			rf.mu.Lock()
		}
	}

	lastEntry := rf.log[len(rf.log)-1]
	voteAvailable := rf.votedFor == -1 || rf.votedFor == args.CandidateId
	upToDate := args.LastLogTerm > lastEntry.Term || (args.LastLogTerm == lastEntry.Term && args.LastLogIndex >= lastEntry.Index)

	if voteAvailable && upToDate {
		rf.votedFor = args.CandidateId

		reply.Term = rf.currentTerm
		reply.VoteGranted = true
	} else {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
	}

	rf.mu.Unlock()

	return nil
}
