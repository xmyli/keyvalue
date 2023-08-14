package keyvalue

type CommandType int

const (
	Get CommandType = iota
	Set
	Append
	Delete
)

type Command struct {
	Type  CommandType
	Done  chan bool
	Key   string
	Value []byte
	TTL   int64
}

type Message struct {
	Data []byte
}
