package handlers

import (
	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/generators"
	"github.com/hop-/pdf-service/internal/kafka"
)

type DocRequest struct {
	Type string         `json:"type"`
	Id   string         `json:"id"`
	Data map[string]any `json:"data"`
}

type DocResponse struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Content string `json:"content"`
}

const (
	ResponseStatusPassed = "passed"
	ResponseStatusFailed = "failed"
)

func sendDocResponse(responsesTopic string, requestId string, status string, content string) error {
	res := DocResponse{
		Id:      requestId,
		Status:  status,
		Content: content,
	}

	msg, err := kafka.NewMessage(&res)
	if err != nil {
		return err
	}

	err = kafka.GetProducer().Send(responsesTopic, msg)
	if err != nil {
		return err
	}

	return nil
}

func NewDocRequestHandler(responsesTopic string) kafka.ConsumerHandle {
	return func(message *kafka.Message) {
		var req DocRequest = DocRequest{}
		err := message.Get(&req)
		golog.Debugf("Kafka request: %+v", req)
		if err != nil {
			golog.Errorf("Failed to parse message: %s", err.Error())
			return
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

			err = sendDocResponse(responsesTopic, req.Id, status, content)
			if err != nil {
				golog.Errorf("Failed to send response via kafka %s", err.Error())
			}

			golog.Infof("Finished report generation for %s request", req.Id)
		}()
	}
}
