package msgqueue

// EventEmitter ...
type EventEmitter interface {
	Emit(event Event) error
}
