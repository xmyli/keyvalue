package main

import (
	"bufio"
	"flag"
	"fmt"
	"keyvalue/internal/keyvalue"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

func main() {
	serversFilePtr := flag.String("f", "", "file that contains address and port of servers")

	flag.Parse()

	file, err := os.Open(*serversFilePtr)
	if err != nil {
		log.Fatalln(err)
	}

	servers := []string{}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		servers = append(servers, scanner.Text())
	}

	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		tokens := strings.Split(input, " ")

		command := strings.ToLower(tokens[0])
		switch command {
		case "get":
			if len(tokens) != 2 {
				fmt.Println("Usage: get [KEY]")
				continue
			}
			key := tokens[1]
			for _, peer := range servers {
				if Get(peer, key) {
					break
				}
			}
		case "set":
			if len(tokens) != 4 {
				fmt.Println("Usage: set [KEY] [VALUE] [TTL]")
				continue
			}
			key := tokens[1]
			value := []byte(tokens[2])
			ttl, err := strconv.ParseInt(tokens[3], 10, 64)
			if err != nil || ttl < 0 {
				fmt.Println("TTL must be a valid number.")
				continue
			}
			for _, peer := range servers {
				if Set(peer, key, value, ttl) {
					break
				}
			}
		case "append":
			if len(tokens) != 4 {
				fmt.Println("Usage: append [KEY] [VALUE] [TTL]")
				continue
			}
			key := tokens[1]
			value := []byte(tokens[2])
			ttl, err := strconv.ParseInt(tokens[3], 10, 64)
			if err != nil || ttl < 0 {
				fmt.Println("TTL must be a valid number.")
				continue
			}
			for _, peer := range servers {
				if Append(peer, key, value, ttl) {
					break
				}
			}
		case "delete":
			if len(tokens) != 2 {
				fmt.Println("Usage: delete [KEY]")
				continue
			}
			key := tokens[1]
			for _, peer := range servers {
				if Delete(peer, key) {
					break
				}
			}
		}
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err().Error())
	}
}

func Get(addr string, key string) bool {
	rpcClient, err := rpc.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return false
	}

	args := keyvalue.GetArgs{}
	args.Key = key

	reply := keyvalue.GetReply{}

	err = rpcClient.Call("KeyValueStore.Get", &args, &reply)
	if err != nil {
		fmt.Println(err)
		rpcClient.Close()
		return false
	}

	if reply.Success {
		fmt.Println("GET:", string(reply.Data))
		rpcClient.Close()
		return true
	}

	rpcClient.Close()
	return false
}

func Set(addr, key string, value []byte, ttl int64) bool {
	rpcClient, err := rpc.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		rpcClient.Close()
		return false
	}

	args := keyvalue.SetArgs{}
	args.Key = key
	args.Value = value
	args.TTL = ttl

	reply := keyvalue.SetReply{}

	err = rpcClient.Call("KeyValueStore.Set", &args, &reply)
	if err != nil {
		fmt.Println(err)
		rpcClient.Close()
		return false
	}

	if reply.Success {
		fmt.Println("SET: SUCCESS")
		rpcClient.Close()
		return true
	}

	rpcClient.Close()
	return false
}

func Append(addr, key string, value []byte, ttl int64) bool {
	rpcClient, err := rpc.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		rpcClient.Close()
		return false
	}

	args := keyvalue.AppendArgs{}
	args.Key = key
	args.Value = value
	args.TTL = ttl

	reply := keyvalue.AppendReply{}

	err = rpcClient.Call("KeyValueStore.Append", &args, &reply)
	if err != nil {
		fmt.Println(err)
		rpcClient.Close()
		return false
	}

	if reply.Success {
		fmt.Println("APPEND: SUCCESS")
		rpcClient.Close()
		return true
	}

	rpcClient.Close()
	return false
}

func Delete(addr, key string) bool {
	rpcClient, err := rpc.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		rpcClient.Close()
		return false
	}

	args := keyvalue.DeleteArgs{}
	args.Key = key

	reply := keyvalue.DeleteReply{}

	err = rpcClient.Call("KeyValueStore.Delete", &args, &reply)
	if err != nil {
		fmt.Println(err)
		rpcClient.Close()
		return false
	}

	if reply.Success {
		fmt.Println("DELETE: SUCCESS")
		rpcClient.Close()
		return true
	}

	rpcClient.Close()
	return false
}
