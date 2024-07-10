package kafka

import (
	"sync"
	"time"

	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/hop-/golog"
)

type SimpleConsumerOptions struct {
	Host    string
	GroupId *string
	Topics  []string
}

type SimpleConsumer struct {
	mu        *sync.Mutex
	consumer  *confluentkafka.Consumer
	isRunning bool
}

func NewSimpleConsumer(options *SimpleConsumerOptions) (*SimpleConsumer, error) {
	config := confluentkafka.ConfigMap{
		"bootstrap.servers": options.Host,
		"auto.offset.reset": "latest",
	}

	if options.GroupId != nil {
		config.SetKey("group.id", *options.GroupId)
	}

	c, err := confluentkafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics(options.Topics, nil)

	return &SimpleConsumer{&sync.Mutex{}, c, true}, nil
}

func (c *SimpleConsumer) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isRunning = false
	c.consumer.Close()
}

func (c *SimpleConsumer) Receive() (*Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	msg, err := c.consumer.ReadMessage(time.Second)
	if err == nil {
		golog.Debugf("New message received on %s", msg.TopicPartition)

		return &Message{msg.Value}, nil
	} else if !err.(confluentkafka.Error).IsTimeout() {
		return nil, err
	}

	return nil, nil
}

func (c *SimpleConsumer) ReceiveUntil() (*Message, error) {
	for c.isRunning {
		msg, err := c.Receive()
		if msg != nil {
			return msg, nil
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
