package lib

// EventCreatedEvent - Published to AMQP
type EventCreatedEvent struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LocationID string `json:""`
}

func (e *EventCreatedEvent) EventName() string {
	return "event.created"
}
