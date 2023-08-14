package keyvalue

import "time"

type SetArgs struct {
	Key   string
	Value []byte
	TTL   int64
}

type SetReply struct {
	Success bool
}

func (kvs *KeyValueStore) Set(args *SetArgs, reply *SetReply) error {
	if !kvs.rf.IsLeader() {
		reply.Success = false
		return nil
	}

	command := Command{}
	command.Type = Set
	command.Done = make(chan bool)
	command.Key = args.Key
	command.Value = args.Value
	command.TTL = args.TTL

	kvs.rf.SendCommand(command)

	<-command.Done

	reply.Success = true
	return nil
}

func (kvs *KeyValueStore) HandleSet(command Command) {
	kvs.Data[command.Key] = command.Value
	kvs.TTL[command.Key] = time.Now().UnixMilli() + command.TTL
	if command.Done != nil {
		command.Done <- true
	}
}
