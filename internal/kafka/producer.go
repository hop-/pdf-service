package kafka

import (
	"sync"

	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerOptions struct {
	Host  string
	Topic string
}

type Producer struct {
	mu       *sync.Mutex
	producer *confluentkafka.Producer
	topic    confluentkafka.TopicPartition
	// TODO: add delivery channel?
}

func NewProducer(options *ProducerOptions) (*Producer, error) {
	config := confluentkafka.ConfigMap{
		"bootstrap.servers": options.Host,
	}

	p, err := confluentkafka.NewProducer(&config)
	if err != nil {
		return nil, err
	}

	topic := confluentkafka.TopicPartition{
		Topic:     &options.Topic,
		Partition: confluentkafka.PartitionAny,
	}

	return &Producer{&sync.Mutex{}, p, topic}, nil
}

func (p *Producer) Send(message *Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.producer.Produce(&confluentkafka.Message{
		TopicPartition: p.topic,
		Value:          message.message,
	}, nil)
}

func (p *Producer) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.producer.Close()
}
