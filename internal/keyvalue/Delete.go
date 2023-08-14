package keyvalue

type DeleteArgs struct {
	Key string
}

type DeleteReply struct {
	Success bool
}

func (kvs *KeyValueStore) Delete(args *DeleteArgs, reply *DeleteReply) error {
	if !kvs.rf.IsLeader() {
		reply.Success = false
		return nil
	}

	command := Command{}
	command.Type = Delete
	command.Done = make(chan bool)
	command.Key = args.Key

	kvs.rf.SendCommand(command)

	<-command.Done

	reply.Success = true
	return nil
}

func (kvs *KeyValueStore) HandleDelete(command Command) {
	delete(kvs.Data, command.Key)
	delete(kvs.TTL, command.Key)
	if command.Done != nil {
		command.Done <- true
	}
}
