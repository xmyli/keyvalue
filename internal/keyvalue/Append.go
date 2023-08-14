package keyvalue

import "time"

type AppendArgs struct {
	Key   string
	Value []byte
	TTL   int64
}

type AppendReply struct {
	Success bool
}

func (kvs *KeyValueStore) Append(args *AppendArgs, reply *AppendReply) error {
	if !kvs.rf.IsLeader() {
		reply.Success = false
		return nil
	}

	command := Command{}
	command.Type = Append
	command.Done = make(chan bool)
	command.Key = args.Key
	command.Value = args.Value
	command.TTL = args.TTL

	kvs.rf.SendCommand(command)

	<-command.Done

	reply.Success = true
	return nil
}

func (kvs *KeyValueStore) HandleAppend(command Command) {
	kvs.Data[command.Key] = append(kvs.Data[command.Key], command.Value...)
	kvs.TTL[command.Key] = time.Now().UnixMilli() + command.TTL
	if command.Done != nil {
		command.Done <- true
	}
}
