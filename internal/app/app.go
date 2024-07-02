package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/generators"
	"github.com/hop-/pdf-service/internal/services"
)

type App struct {
	exitChan  chan os.Signal
	options   Options
	services  []services.Service
	generator *generators.ConcurrentPdfGenerator
	wg        *sync.WaitGroup
}

// App constructor
func NewApp(options Options) *App {
	golog.Debugf("App options are: %+v", options)

	generator := generators.NewConcurrentPdfGenerator(options.Concurrency, options.EngineType)

	srvs := []services.Service{}

	if options.Http.Enabled {
		srvs = append(srvs, services.NewHttpService(
			options.Http.Port,
			options.Http.Secure,
			options.Http.Cert,
			options.Http.Key,
			generator,
		))
	}

	if options.Kafka.Enabled {
		srvs = append(srvs, services.NewKafkaService(
			options.Kafka.Host,
			options.Kafka.ConsumerGroupId,
			options.Kafka.RequestsTopic,
			options.Kafka.ResponsesTopic,
			generator,
		))
	}

	app := App{
		make(chan os.Signal, 1),
		options,
		srvs,
		generator,
		new(sync.WaitGroup),
	}

	// Signal handling
	signal.Notify(app.exitChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return &app
}

func New(optionModifiers ...OptionModifier) *App {
	o := defaultOptions()
	for _, omd := range optionModifiers {
		omd(&o)
	}

	return NewApp(o)
}

// App start function
func (a *App) Start() {
	// Graceful shutdown
	go a.gracefulShutDownTracker()

	for _, s := range a.services {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			// This is available Go versions >= 1.22
			s.Start()
		}()
	}

	// Wait until all goroutines are done
	a.wg.Wait()
}

func (a *App) gracefulShutDownTracker() {
	<-a.exitChan

	// Iterate and stop all services
	for _, s := range a.services {
		s.Stop()
	}
}
