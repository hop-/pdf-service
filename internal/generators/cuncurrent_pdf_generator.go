package generators

import (
	"sync"

	"github.com/hop-/pdf-service/internal/reports"
)

type ConcurrentPdfGenerator struct {
	concurrencyMutex chan bool
	engine           string
}

var cpgLock = &sync.Mutex{}
var cpgInstance *ConcurrentPdfGenerator

func GetConcurrentPdfGenerator() *ConcurrentPdfGenerator {
	if cpgInstance != nil {
		return cpgInstance
	}

	return InitCpgInstnace(1, "chromedp")
}

func InitCpgInstnace(concurrency uint8, engine string) *ConcurrentPdfGenerator {
	cpgLock.Lock()
	defer cpgLock.Unlock()

	cpgInstance = &ConcurrentPdfGenerator{
		concurrencyMutex: make(chan bool, concurrency),
		engine:           engine,
	}

	return cpgInstance
}

func (g *ConcurrentPdfGenerator) Generate(template string, data map[string]any) ([]byte, error) {
	// Concurrent mutex lock
	g.concurrencyMutex <- true

	// Unlock
	defer func() {
		<-g.concurrencyMutex
	}()

	reportGenerator, err := reports.NewReportGenerator(template, g.engine)
	if err != nil {
		return []byte{}, err
	}

	content, err := reportGenerator.Generate(data)
	if err != nil {
		return []byte{}, err
	}

	return content, err
}

func (g *ConcurrentPdfGenerator) GenerateBase64(template string, data map[string]any) (string, error) {
	// Concurrent mutex lock
	g.concurrencyMutex <- true

	// Unlock
	defer func() {
		<-g.concurrencyMutex
	}()

	reportGenerator, err := reports.NewReportGenerator(template, g.engine)
	if err != nil {
		return "", err
	}

	contentBase64, err := reportGenerator.GenerateBase64(data)
	if err != nil {
		return "", err
	}

	return contentBase64, err
}
