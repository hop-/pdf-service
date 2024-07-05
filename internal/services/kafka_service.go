package services

import (
	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/handlers"
	"github.com/hop-/pdf-service/internal/kafka"
)

type KafkaService struct {
	consumer *kafka.Consumer
}

func NewKafkaService(
	host string,
	consumerGroupId string,
	requestsTopic string,
	responsesTopic string,
) *KafkaService {
	// Create Consumer
	consumerOptions := kafka.ConsumerOptions{
		Host:    host,
		GroupId: &consumerGroupId,
		Topics:  []string{requestsTopic},
		Handler: handlers.NewDocRequestHandler(responsesTopic),
	}

	consumer, err := kafka.NewConsumer(&consumerOptions)
	if err != nil {
		golog.Fatalf("Failed to create kafka consumer: %s", err.Error())
	}

	err = kafka.InitProducerOnce(host)
	if err != nil {
		golog.Fatalf("Failed to create kafka producer: %s", err.Error())
	}

	return &KafkaService{
		consumer: consumer,
	}
}

func (s *KafkaService) Start() {
	golog.Info("Listening for Kafka messages")

	// Message gethering loop
	s.consumer.Run()
}

func (s *KafkaService) Stop() {
	// Close consumer on exit
	if s.consumer != nil {
		s.consumer.Close()
		s.consumer = nil
	}
	// Close producer on exit
	kafka.GetProducer().Close()
}
