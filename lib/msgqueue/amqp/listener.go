package amqp

import (
	"bookmyevents/lib"
	"bookmyevents/lib/events"
	"bookmyevents/lib/msgqueue"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type amqpEventListener struct {
	connection *amqp.Connection
	queue      string
}

func (a *amqpEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)

	return err
}

func NewAMQPEventListener(conn *amqp.Connection, queue string) (msgqueue.EventListener, error) {
	listener := &amqpEventListener{
		connection: conn,
		queue:      queue,
	}

	err := listener.setup()
	if err != nil {
		return nil, err
	}

	return listener, err
}

func (a *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	defer channel.Close()

	for _, eventName := range eventNames {
		if err := channel.QueueBind(a.queue, eventName, "events", false, nil); err != nil {
			return nil, nil, err
		}
	}

	msgs, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	eventsChan := make(chan msgqueue.Event)
	errorsChan := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errorsChan <- fmt.Errorf("msg did not contain x-event-name header")
				msg.Nack(false, true)
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errorsChan <- fmt.Errorf("x-event-name header is not a string, but %t", rawEventName)
				msg.Nack(false, true)
				continue
			}

			var event msgqueue.Event

			switch eventName {
			case lib.EVENTCREATED:
				event = new(events.EventCreatedEvent)
			default:
				errorsChan <- fmt.Errorf("Event type %s is unknown", eventName)
				continue
			}
			err := json.Unmarshal(msg.Body, &event)
			if err != nil {
				log.Println("Error in unmarshalling queue message -" + err.Error())
				errorsChan <- err
				continue
			}

			eventsChan <- event

		}

	}()

	return eventsChan, errorsChan, nil
}
