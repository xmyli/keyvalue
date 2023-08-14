package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"keyvalue/internal/keyvalue"
	"keyvalue/pkg/raft"
	"log"
	"os"
)

func main() {
	gob.Register(keyvalue.Command{})

	peersFilePtr := flag.String("f", "", "file that contains address and port of peers")
	indexPtr := flag.Int("i", -1, "index of this peer in the peers file")
	raftPortPtr := flag.String("rp", "8000", "port of raft server")
	storePortPtr := flag.String("sp", "8001", "port of key-value store")

	flag.Parse()

	file, err := os.Open(*peersFilePtr)
	if err != nil {
		log.Fatalln(err)
	}

	peers := []string{}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		peers = append(peers, scanner.Text())
	}

	commandsCh := make(chan interface{})

	rf := raft.CreateRaftServer(*raftPortPtr, peers, *indexPtr, commandsCh)
	rf.Start()

	kvs := keyvalue.CreateKeyValueStore(*storePortPtr, rf, commandsCh)
	kvs.Start()

	select {}
}
