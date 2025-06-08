package order

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type EventData struct {
	AccountId string `json:"user_id"`
	ProductId string `json:"product_id"`
}

type Event struct {
	Type      string    `json:"type"`
	EventData EventData `json:"data"`
}

var done = make(chan bool)

func (o orderService) SendMessageToRecommender(event Event, topic string) error {
	log.Printf("SendMessageToRecommender called for topic: %s", topic)

	jsonMessage, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshalling event: %v", err)
		return err
	}

	log.Printf("Marshalled event to JSON: %s", string(jsonMessage))

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMessage),
	}

	log.Printf("Sending message to Kafka topic %s...", topic)
	o.producer.Input() <- msg
	log.Printf("Message sent to producer input channel")

	return nil
}

func (o orderService) MsgHandler() {
	go func() {
		for {
			select {
			case success := <-o.producer.Successes():
				log.Printf("Message sent to %s successfully: %v", success.Topic, success.Partition)
			case err := <-o.producer.Errors():
				log.Printf("Error sending message to %s: %v", err.Msg.Topic, err.Err)
			case <-done:
				log.Printf("Producer closed for %s", "order_events")
				return
			}
		}
	}()
}

func (o orderService) Close() {
	if err := o.producer.Close(); err != nil {
		log.Printf("Error closing producer: %v", err)
	} else {
		done <- true
	}
}
