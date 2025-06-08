package product

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type EventData struct {
	ID          *string  `json:"product_id"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	AccountID   *string  `json:"accountID"`
}

type Event struct {
	Type string    `json:"type"`
	Data EventData `json:"data"`
}

var done = make(chan bool)

func (p productService) SendMessageToRecommender(event Event, topic string) error {
	jsonMessage, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshalling event: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMessage),
	}

	p.producer.Input() <- msg

	return nil
}

func (p productService) MsgHandler() {
	go func() {
		for {
			select {
			case success := <-p.producer.Successes():
				log.Printf("Message sent to %s successfully: %v", success.Topic, success.Partition)
			case err := <-p.producer.Errors():
				log.Printf("Error sending message to %s: %v", err.Msg.Topic, err.Err)
			case <-done:
				log.Printf("Producer closed for %s", "product_events")
				return
			}
		}
	}()
}

func (p productService) Close() {
	if err := p.producer.Close(); err != nil {
		log.Printf("Error closing producer: %v", err)
	} else {
		done <- true
	}
}
