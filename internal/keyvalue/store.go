package keyvalue

import (
	"keyvalue/pkg/raft"
	"log"
	"net"
	"net/rpc"
)

type KeyValueStore struct {
	Data map[string][]byte
	TTL  map[string]int64

	port string
	rf   *raft.Raft
	ch   chan interface{}
}

func (kvs *KeyValueStore) HandleCommands() {
	for {
		commandInterface := <-kvs.ch
		command := commandInterface.(Command)
		switch command.Type {
		case Get:
			continue
		case Set:
			kvs.HandleSet(command)
		case Append:
			kvs.HandleAppend(command)
		case Delete:
			kvs.HandleDelete(command)
		}
	}
}

func (kvs *KeyValueStore) Start() {
	rpcs := rpc.NewServer()
	rpcs.Register(kvs)
	listener, err := net.Listen("tcp", ":"+kvs.port)
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

	go kvs.HandleCommands()
}

func CreateKeyValueStore(port string, rf *raft.Raft, ch chan interface{}) *KeyValueStore {
	keyValueStore := KeyValueStore{}

	keyValueStore.Data = map[string][]byte{}
	keyValueStore.TTL = map[string]int64{}

	keyValueStore.port = port
	keyValueStore.rf = rf
	keyValueStore.ch = ch

	return &keyValueStore
}
