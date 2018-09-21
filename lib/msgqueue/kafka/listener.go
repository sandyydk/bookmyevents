package kafka

import (
	evts "bookmyevents/lib/events"
	"bookmyevents/lib/msgqueue"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/Shopify/sarama"
)

type kafkaEventListener struct {
	consumer   sarama.Consumer
	partitions []int32
}

func NewKafkaEventListener(client sarama.Client, partitions []int32) (msgqueue.EventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	listener := &kafkaEventListener{
		consumer:   consumer,
		partitions: partitions,
	}
	return listener, nil
}

func (e *kafkaEventListener) Listen(events ...string) (<-chan msgqueue.Event, <-chan error, error) {
	var err error

	topic := "events"
	results := make(chan msgqueue.Event)
	errors := make(chan error)

	partitions := e.partitions

	if len(partitions) == 0 {
		partitions, err = e.consumer.Partitions(topic)
		if err != nil {
			return nil, nil, err
		}
	}

	log.Printf("topic %s has partitions: %v", topic, partitions)

	for _, partition := range partitions {
		con, err := e.consumer.ConsumePartition(topic, partition, 0)
		if err != nil {
			return nil, nil, err
		}

		go func() {
			for msg := range con.Messages() {
				body := messageEnvelope{}
				err := json.Unmarshal(msg.Value, &body)
				if err != nil {
					errors <- fmt.Errorf("could not JSON-decode message: %s", err)
					continue
				}

				var event msgqueue.Event
				switch body.EventName {
				case "event.created":
					event = &evts.EventCreatedEvent{}
				case "location.created":
					event = &evts.LocationCreatedEvent{}
				default:
					errors <- fmt.Errorf("unknown event type: %s", body.EventName)
					continue
				}

				cfg := mapstructure.DecoderConfig{
					Result:  event,
					TagName: "json",
				}

				decoder, err := mapstructure.NewDecoder(&cfg)
				if err != nil {
					errors <- fmt.Errorf("could not map event %s: %s", body.EventName, err)
				}
				err = decoder.Decode(body.Payload)
				if err != nil {
					errors <- fmt.Errorf("could not map event %s: %s", body.EventName, err)
				}
			}
		}()
	}
	return results, errors, nil
}
