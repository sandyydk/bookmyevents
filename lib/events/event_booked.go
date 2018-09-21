package events

type EventBookedEvent struct {
	EventID string `json:"eventId"`
	UserID  string `json:"userId"`
}

func (c *EventBookedEvent) EventName() string {
	return "event.booked"
}

func (c *EventBookedEvent) Partition() string {
	return c.EventID
}
