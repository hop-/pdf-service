package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/kafka"
	"github.com/hop-/pdf-service/internal/reports"
)

type HttpOptions struct {
	Enabled bool
	Port    int
	Secure  bool
	Cert    string
	Key     string
}

type KafkaOptions struct {
	Enabled         bool
	Host            string
	ConsumerGroupId string
	RequestsTopic   string
	ResponsesTopic  string
}

type Options struct {
	EngineType  string
	Concurrency uint8
	Http        HttpOptions
	Kafka       KafkaOptions
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
	wg               *sync.WaitGroup
}

// App constructor
func NewApp(options Options) App {
	golog.Debugf("App options are: %+v", options)

	o := App{
		make(chan os.Signal, 1),
		make(chan bool, options.Concurrency),
		options,
		[]ShutdownHandlerFunc{},
		nil,
		nil,
		true,
		new(sync.WaitGroup),
	}

	// Signal handling
	signal.Notify(o.exitChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return o
}

// App start function
func (a *App) Start() {
	// Graceful shutdown
	go a.gracefulShutDownTracker()

	if a.options.Kafka.Enabled {
		a.wg.Add(1)
		go a.startKafka()
	}
	if a.options.Http.Enabled {
		a.wg.Add(1)
		go a.startHttp()
	}

	// Wait until all goroutines are done
	a.wg.Wait()
}

func (a *App) generateReport(req KafkaRequest, engine string) (string, error) {
	// Concurrent mutex lock
	a.concurrencyMutex <- true

	// Unlock
	defer func() {
		<-a.concurrencyMutex
	}()

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

func (a *App) OnShutdown(h ShutdownHandlerFunc) {
	a.shutdownHandlers = append(a.shutdownHandlers, h)
}

func (a *App) gracefulShutDownTracker() {
	<-a.exitChan
	a.isRunning = false

	// Iterate and run shutdown handlers
	for i := range a.shutdownHandlers {
		a.shutdownHandlers[i]()
	}
}
