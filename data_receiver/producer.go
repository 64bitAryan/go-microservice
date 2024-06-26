package main

import (
	"encoding/json"
	"fmt"

	"github.com/64bitAryan/go-microservice/types"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type DataProducer interface {
	ProduceData(types.OBUDATA) error
}

type Kafkaproducer struct {
	producer *kafka.Producer
}

func NewkafkaProducer() (*Kafkaproducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// start anoter go-routine to check if we have delivered the data
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &Kafkaproducer{
		producer: p,
	}, nil
}

func (p *Kafkaproducer) ProduceData(data types.OBUDATA) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaTopic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)
}
