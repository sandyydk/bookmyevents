package msgqueue

type Event interface {
	EventName() string
	PartitionKey() string
}
