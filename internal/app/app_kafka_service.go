package app

import (
	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/kafka"
)

func (a *App) startKafka() {
	defer a.wg.Done()
	// Create Consumer
	consumerOptions := kafka.ConsumerOptions{
		Host:    a.options.Kafka.Host,
		GroupId: &a.options.Kafka.ConsumerGroupId,
		Topics:  []string{a.options.Kafka.RequestsTopic},
	}
	consumer, err := kafka.NewConsumer(&consumerOptions)
	if err != nil {
		golog.Fatalf("Failed to create kafka consumer: %s", err.Error())
	}
	// Close consumer on exit
	a.OnShutdown(func() {
		consumer.Close()
	})
	a.consumer = consumer

	// Create Producer
	producerOptions := kafka.ProducerOptions{
		Host:  a.options.Kafka.Host,
		Topic: a.options.Kafka.ResponsesTopic,
	}
	prodcuer, err := kafka.NewProducer(&producerOptions)
	if err != nil {
		golog.Fatalf("Failed to create kafka producer: %s", err.Error())
	}
	a.OnShutdown(func() {
		prodcuer.Close()
	})
	a.producer = prodcuer

	golog.Info("Listening for Kafka messages")
	// Message gethering loop
	for a.isRunning {
		message, err := consumer.ReceiveUntil()
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
			content, err := a.generateReport(req, a.options.EngineType)
			if err != nil {
				status = ResponseStatusFailed
				golog.Errorf("Failed to generate report: %s", err)
			}

			err = a.sendKafkaResponse(req.Id, status, content)
			if err != nil {
				golog.Errorf("Failed to send response via kafka %s", err.Error())
			}

			golog.Infof("Finished report generation for %s request", req.Id)
		}()
	}
}

func (a *App) sendKafkaResponse(requestId string, status string, content string) error {
	res := KafkaResponse{
		Id:      requestId,
		Status:  status,
		Content: content,
	}

	msg, err := kafka.NewMessage(&res)
	if err != nil {
		return err
	}

	err = a.producer.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
