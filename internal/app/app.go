package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/kafka"
	"github.com/hop-/pdf-service/internal/reports"
)

type KafkaOptions struct {
	Host            string
	ConsumerGroupId string
	RequestsTopic   string
	ResponsesTopic  string
}
type Options struct {
	EngineType     string
	Concurrency    uint8
	KafkaIsEnabled bool
	Kafka          KafkaOptions
}

type ShutdownHandlerFunc = func()

type App struct {
	exitChan         chan os.Signal
	concurrencyMutex chan bool
	options          Options
	shutdownHandlers []ShutdownHandlerFunc
	consumer         *kafka.Consumer
	producer         *kafka.Producer
	isRunning        bool
}

// App constructor
func NewApp(options Options) App {
	o := App{
		make(chan os.Signal, 1),
		make(chan bool, options.Concurrency),
		options,
		[]ShutdownHandlerFunc{},
		nil,
		nil,
		true,
	}

	// Signal handling
	signal.Notify(o.exitChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return o
}

// App start function
func (a *App) Start() {
	// Graceful shutdown
	go a.shutDown()

	// TODO: use kafka is enabled flag
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

	// Message gethering loop
	for a.isRunning {
		message, err := consumer.ReceiveUntil()
		if err != nil {
			golog.Errorf("Failed to read message: %s", err.Error())
			continue
		} else if message == nil {
			golog.Errorf("message is empty or nil")
			continue
		}

		var req ReportRequest = ReportRequest{}
		err = message.Get(&req)
		golog.Debug(req)
		if err != nil {
			golog.Errorf("Failed to parse message: %s", err.Error())
			continue
		}

		golog.Infof("New report requested for %s with %s id", req.Type, req.Id)

		// Generate report concurrently
		go func() {
			// Concurrent mutex lock
			a.concurrencyMutex <- true

			golog.Infof("Starting report generation for %s request", req.Id)

			// Generate report
			status := ResponseStatusPassed
			content, err := a.generateReport(req, a.options.EngineType)
			if err != nil {
				status = ResponseStatusFailed
				golog.Errorf("Failed to generate report: %s", err)
			}

			err = a.sendResponse(req.Id, status, content)
			if err != nil {
				golog.Errorf("Failed to send response via kafka %s", err.Error())
			}

			golog.Infof("Finished report generation for %s request", req.Id)

			// Unlock
			<-a.concurrencyMutex
		}()
	}
}

func (a *App) generateReport(req ReportRequest, engine string) (string, error) {
	data, err := req.JsonData()
	if err != nil {
		return "", err
	}

	reportGenerator, err := reports.NewReportGenerator(req.Type, engine)
	if err != nil {
		return "", err
	}

	contentBase64, err := reportGenerator.GenerateBase64(data)
	if err != nil {
		return "", err
	}

	return contentBase64, err
}

func (a *App) sendResponse(requestId string, status string, content string) error {
	res := ReportResponse{
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

func (a *App) OnShutdown(h ShutdownHandlerFunc) {
	a.shutdownHandlers = append(a.shutdownHandlers, h)
}

func (a *App) shutDown() {
	<-a.exitChan
	a.isRunning = false

	// Iterate and run shutdown handlers
	for i := range a.shutdownHandlers {
		a.shutdownHandlers[i]()
	}
}
