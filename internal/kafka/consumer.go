package kafka

import (
	"time"

	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/hop-/golog"
)

type ConsumerOptions struct {
	Host    string
	GroupId *string
	Topics  []string
}

type Consumer struct {
	consumer  *confluentkafka.Consumer
	isRunning bool
}

func NewConsumer(options *ConsumerOptions) (*Consumer, error) {
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

	return &Consumer{c, true}, nil
}

func (c *Consumer) Close() {
	c.isRunning = false
	c.consumer.Close()
}

func (c *Consumer) Receive() (*Message, error) {
	msg, err := c.consumer.ReadMessage(time.Second)
	if err == nil {
		golog.Debugf("New message received on %s", msg.TopicPartition)

		return &Message{msg.Value}, nil
	} else if !err.(confluentkafka.Error).IsTimeout() {
		return nil, err
	}

	return nil, nil
}

func (c *Consumer) ReceiveUntil() (*Message, error) {
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
