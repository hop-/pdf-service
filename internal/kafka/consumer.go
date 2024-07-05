package kafka

import (
	"github.com/hop-/golog"
)

type ConsumerHandle func(*Message)

type ConsumerOptions struct {
	Host    string
	GroupId *string
	Topics  []string
	Handler ConsumerHandle
}

type Consumer struct {
	consumer  *SimpleConsumer
	handler   ConsumerHandle
	isRunning bool
}

func NewConsumer(opts *ConsumerOptions) (*Consumer, error) {
	simpleOpts := SimpleConsumerOptions{
		Host:    opts.Host,
		GroupId: opts.GroupId,
		Topics:  opts.Topics,
	}

	consumer, err := NewSimpleConsumer(&simpleOpts)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:  consumer,
		handler:   opts.Handler,
		isRunning: false,
	}, nil
}

func (c *Consumer) Run() {
	c.isRunning = true
	// Message gethering loop
	for c.isRunning {
		message, err := c.consumer.ReceiveUntil()
		if err != nil {
			golog.Errorf("Failed to read message: %s", err.Error())
			continue
		} else if message == nil {
			golog.Debug("Message is empty or nil")
			continue
		}

		c.handler(message)
	}
}

func (c *Consumer) Close() {
	c.isRunning = false
	c.consumer.Close()
}
