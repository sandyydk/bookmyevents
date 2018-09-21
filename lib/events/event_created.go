package events

import (
	"bookmyevents/lib"
	"time"
)

// EventCreatedEvent - Published to AMQP
type EventCreatedEvent struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	LocationID string    `json:"location_id"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
}

func (e *EventCreatedEvent) EventName() string {
	return lib.EVENTCREATED
}

func (e *EventCreatedEvent) PartitionKey() string {
	return e.ID
}
