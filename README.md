# Distributed Key-Value Store

Distributed key-value store that uses the Raft consensus algorithm for fault-tolerance. Supports time-to-live (TTL) for key-value pairs.

## Usage
Requires Go to be installed.
1. Build the server.
   ```
   go build ./cmd/keyvalue
   ```
2. Create a file that defines the address and port for each peer.
   ```
   127.0.0.1:8000
   127.0.0.1:8001
   127.0.0.1:8002
   ```
3. Start the key-value store server. Requires at least 2 instances to be accessible. 2n+1 instances is required to tolerate a maximum of n faults.
   ```
   ./keyvalue -f=FILENAME -i=NUMBER -rp=PORT -sp=PORT
   ```
   -f=FILENAME: The file that defines the address and port for each peer (Step 2)
   -i=NUMBER: The index of this peer in the above file. Starts on index 0.
   -rp=PORT: The port of the Raft server, should be the same as the port defined in the above file. Only for internal communications.
   -sp=PORT: The port of the key-value store server. This will be used by external clients.

## Testing
Includes a CLI tool for manual testing. Uses pre-defined addresses in peers.conf and servers.conf.
1. Build the CLI tool.
    ```
    go build ./cmd/test
    ```
2. Start 5 peers using the addresses in peers.conf and servers.conf.
    ```
    ./keyvalue -f=peers.conf -i=0 -rp=8000 -sp=3000 &
    ./keyvalue -f=peers.conf -i=0 -rp=8001 -sp=3001 &
    ./keyvalue -f=peers.conf -i=0 -rp=8002 -sp=3002 &
    ./keyvalue -f=peers.conf -i=0 -rp=8003 -sp=3003 &
    ./keyvalue -f=peers.conf -i=0 -rp=8004 -sp=3004 &
    ```
3. Start the CLI tool. Refer to the commands section for available commands.
    ```
    ./test -f=servers.conf
    ```

## Commands
4 different commands are supported, accessed through either RPC or commands from the CLI tool.
args and reply structs should be imported from "keyvalue/internal/keyvalue".
Each call uses a different set of args and reply structs.



### get
Gets the value of the specified key, if not expired.

#### RPC
```
.Call("KeyValueStore.Get", &args, &reply)
```
args is of type keyvalue.GetArgs and reply is of type keyvalue.GetReply

##### Command
```
get [KEY]
```



### set
Sets a key-value pair, requires TTL to be specified.

##### RPC
```
.Call("KeyValueStore.Set", &args, &reply)
```
args is of type keyvalue.SetArgs and reply is of type keyvalue.SetReply

##### Command
```
set [KEY] [VALUE] [TTL]
```



### append
Appends the new value to the end of an existing value. Requires TTL to be specified, which replaces the current expiry time. Use an empty new value to refresh the value by TTL.

##### RPC
```
.Call("KeyValueStore.Append", &args, &reply)
```
args is of type keyvalue.AppendArgs and reply is of type keyvalue.AppendReply

##### Command
```
append [KEY] [VALUE] [TTL]
```



### delete
Deletes a key-value pair.

##### RPC
```
.Call("KeyValueStore.Delete", &args, &reply)
```
args is of type keyvalue.DeleteArgs and reply is of type keyvalue.DeleteReply

##### Command
```
delete [KEY]
```
