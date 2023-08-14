package keyvalue

import "time"

type GetArgs struct {
	Key string
}

type GetReply struct {
	Success bool
	Data    []byte
}

func (kvs *KeyValueStore) Get(args *GetArgs, reply *GetReply) error {
	if !kvs.rf.IsLeader() {
		return nil
	}

	data, exists := kvs.Data[args.Key]
	if !exists {
		return nil
	}

	if kvs.TTL[args.Key] < time.Now().UnixMilli() {
		return nil
	}

	reply.Success = true
	reply.Data = data
	return nil
}
