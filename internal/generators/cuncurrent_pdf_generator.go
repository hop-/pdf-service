package generators

import "github.com/hop-/pdf-service/internal/reports"

type ConcurrentPdfGenerator struct {
	concurrencyMutex chan bool
	engine           string
}

func NewConcurrentPdfGenerator(concurrency uint8, engine string) *ConcurrentPdfGenerator {
	return &ConcurrentPdfGenerator{
		concurrencyMutex: make(chan bool, concurrency),
		engine:           engine,
	}
}

func (g *ConcurrentPdfGenerator) Generate(template string, data map[string]any) (string, error) {
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
