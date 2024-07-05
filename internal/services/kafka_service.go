package services

import (
	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/generators"
	"github.com/hop-/pdf-service/internal/kafka"
)

type KafkaRequest struct {
	Type string         `json:"type"`
	Id   string         `json:"id"`
	Data map[string]any `json:"data"`
}

type KafkaResponse struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Content string `json:"content"`
}

const (
	ResponseStatusPassed = "passed"
	ResponseStatusFailed = "failed"
)

type KafkaService struct {
	consumer  *kafka.Consumer
	producer  *kafka.Producer
	isRunning bool
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
	}

	consumer, err := kafka.NewConsumer(&consumerOptions)
	if err != nil {
		golog.Fatalf("Failed to create kafka consumer: %s", err.Error())
	}

	// Create Producer
	producerOptions := kafka.ProducerOptions{
		Host:  host,
		Topic: responsesTopic,
	}

	prodcuer, err := kafka.NewProducer(&producerOptions)
	if err != nil {
		golog.Fatalf("Failed to create kafka producer: %s", err.Error())
	}

	return &KafkaService{
		consumer:  consumer,
		producer:  prodcuer,
		isRunning: false,
	}
}

func (s *KafkaService) Start() {
	s.isRunning = true
	golog.Info("Listening for Kafka messages")
	// Message gethering loop
	for s.isRunning {
		message, err := s.consumer.ReceiveUntil()
		if err != nil {
			golog.Errorf("Failed to read message: %s", err.Error())
			continue
		} else if message == nil {
			golog.Debug("Message is empty or nil")
			continue
		}

		var req KafkaRequest = KafkaRequest{}
		err = message.Get(&req)
		golog.Debugf("Kafka request: %+v", req)
		if err != nil {
			golog.Errorf("Failed to parse message: %s", err.Error())
			continue
		}

		golog.Infof("New report requested for %s with %s id", req.Type, req.Id)

		// Generate report concurrently
		go func() {
			golog.Infof("Starting report generation for %s request", req.Id)

			// Generate report
			status := ResponseStatusPassed

			content, err := generators.GetConcurrentPdfGenerator().Generate(req.Type, req.Data)
			if err != nil {
				status = ResponseStatusFailed
				golog.Errorf("Failed to generate report: %s", err)
			}

			err = s.send(req.Id, status, content)
			if err != nil {
				golog.Errorf("Failed to send response via kafka %s", err.Error())
			}

			golog.Infof("Finished report generation for %s request", req.Id)
		}()
	}
}

func (s *KafkaService) Stop() {
	s.isRunning = false
	// Close consumer on exit
	if s.consumer != nil {
		s.consumer.Close()
		s.consumer = nil
	}
	// Close producer on exit
	if s.producer != nil {
		s.producer.Close()
		s.producer = nil
	}
}

func (s *KafkaService) send(requestId string, status string, content string) error {
	res := KafkaResponse{
		Id:      requestId,
		Status:  status,
		Content: content,
	}

	msg, err := kafka.NewMessage(&res)
	if err != nil {
		return err
	}

	err = s.producer.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
