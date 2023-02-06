package consumer

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"survey-ussd/core/environment"
	"survey-ussd/core/topics"
	"time"
)

type Consumer struct {
	client  *kafka.Consumer
	env     *environment.Environment
	topics  *topics.Topics
	context context.Context
}

func logf(msg string, a ...interface{}) {
	go func() {
		fmt.Printf(msg, a...)
		println()
	}()
}

func New(env *environment.Environment, context context.Context, topics *topics.Topics) *Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": env.KafkaBroker,
		"group.id":          env.KafkaGroupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}
	importedTopics := make([]string, 0, len(topics.List))
	for k := range topics.List {
		importedTopics = append(importedTopics, k)
	}

	CreateTopics(importedTopics, env.KafkaBroker)
	c.SubscribeTopics(importedTopics, nil)
	return &Consumer{env: env, client: c, topics: topics, context: context}
}

func (c *Consumer) Consume() {
	for true {
		m, err := c.client.ReadMessage(time.Second)
		if err == nil {
			go c.topics.List[*m.TopicPartition.Topic](m.Value, *m.TopicPartition.Topic)
		} else if err.(kafka.Error).Error() != kafka.ErrTimedOut.String() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.

			fmt.Printf("Consumer error: %v (%v)\n", err, m)
		}
	}
	c.client.Close()
}

func CreateTopics(topics []string, broker string) {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		panic(err)
	}

	// Create topics
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		panic("ParseDuration(60s)")
	}

	for _, topic := range topics {
		results, err := adminClient.CreateTopics(context.Background(),
			// Multiple topics can be created simultaneously
			// by providing more TopicSpecification structs here.
			[]kafka.TopicSpecification{{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1}},
			// Admin options
			kafka.SetAdminOperationTimeout(maxDur))

		if err != nil {
			fmt.Printf("Failed to create topic: %v\n", err)
		} else {
			for _, result := range results {
				fmt.Printf("%s\n", result)
			}
		}
	}

	adminClient.Close()
}
