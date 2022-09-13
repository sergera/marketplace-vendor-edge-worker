package evt

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sergera/marketplace-vendor-edge-worker/internal/conf"
	"github.com/sergera/marketplace-vendor-edge-worker/internal/domain"
)

var Topics map[domain.Status]string = map[domain.Status]string{
	domain.Unconfirmed: "orders__unconfirmed",
	domain.InProgress:  "orders__in_progress",
	domain.Ready:       "orders__ready",
	domain.InTransit:   "orders__in_transit",
	domain.Delivered:   "orders__delivered",
}

type EventHandler struct {
	consumer *kafka.Consumer
	producer *kafka.Producer
}

func NewEventHandler() *EventHandler {
	conf := conf.GetConf()
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": conf.KafkaHost + ":" + conf.KafkaPort,
	})
	if err != nil {
		log.Panic(err)
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": conf.KafkaHost + ":" + conf.KafkaPort,
		"group.id":          "marketplace-vendor-edge-worker",
	})
	if err != nil {
		log.Panic(err)
	}

	eventHandler := EventHandler{c, p}
	go eventHandler.reportDeliveries()

	return &eventHandler
}

func (e *EventHandler) CloseProducer() {
	e.producer.Close()
}

func (e *EventHandler) CloseConsumer() {
	e.consumer.Close()
}

func (e *EventHandler) Flush(timeoutMs int) {
	// Wait for message deliveries before shutting down
	e.producer.Flush(timeoutMs)
}

func (e *EventHandler) reportDeliveries() {
	// Delivery report handler for produced messages
	for e := range e.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}

func (e *EventHandler) Produce(topic string, key string, value []byte) {
	e.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	},
		nil)
}

func (e *EventHandler) Consume(topics []string, callback func(kafka.Message)) {
	e.consumer.SubscribeTopics(topics, nil)
	for {
		msg, err := e.consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			callback(*msg)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("consumer error: %v (%v)\n", err, msg)
		}
	}
}
