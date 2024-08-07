package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/64bitAryan/go-microservice/aggregator/client"
	"github.com/64bitAryan/go-microservice/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
	// creating consumer, **consuming on localhost**
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	// A signal handler or similar could be used to set this to false to break the loop.
	return &KafkaConsumer{
		consumer:    c,
		isRunning:   true,
		calcService: svc,
		aggClient:   aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("Kafka transport started...")
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("Kafka consumer error %s ", err)
			continue
		}
		var data types.OBUDATA
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON searialization error: %s", err)
			continue
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error: %s", err)
			continue
		}
		req := &types.AggregateRequest{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			ObuId: int32(data.OBUID),
		}
		if err := c.aggClient.Aggregate(context.Background(), req); err != nil {
			logrus.Errorf("aggregate error: %s", err)
			continue
		}
	}
}
