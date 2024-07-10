package kafka

import (
	"fmt"
	"sync"

	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type GenericProducer struct {
	mu       *sync.Mutex
	producer *confluentkafka.Producer
}

var producerLock = &sync.Mutex{}
var producerInstance *GenericProducer

func InitProducerOnce(host string) error {
	if producerInstance != nil {
		return fmt.Errorf("Producer has already been initialized")
	}

	producerLock.Lock()
	defer producerLock.Unlock()

	config := confluentkafka.ConfigMap{
		"bootstrap.servers": host,
	}

	p, err := confluentkafka.NewProducer(&config)
	if err != nil {
		return err
	}

	producerInstance = &GenericProducer{&sync.Mutex{}, p}

	return nil
}

func GetProducer() *GenericProducer {
	if producerInstance == nil {
		panic("Producer should be initialized first")
	}

	return producerInstance
}

func (p *GenericProducer) Send(topic string, message *Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	topicPartition := confluentkafka.TopicPartition{
		Topic:     &topic,
		Partition: confluentkafka.PartitionAny,
	}

	return p.producer.Produce(&confluentkafka.Message{
		TopicPartition: topicPartition,
		Value:          message.message,
	}, nil)
}

func (p *GenericProducer) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.producer.Close()
}
