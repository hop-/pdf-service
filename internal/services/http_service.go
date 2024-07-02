package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/generators"
)

type HttpService struct {
	srv       *http.Server
	tls       bool
	certFile  string
	keyFile   string
	generator *generators.ConcurrentPdfGenerator
}

func NewHttpService(
	port uint16,
	tls bool,
	certFile string,
	keyFile string,
	generator *generators.ConcurrentPdfGenerator,
) *HttpService {
	addr := fmt.Sprintf(":%d", port)

	var router http.Handler
	// TODO: add router

	srv := http.Server{Addr: addr, Handler: router}

	return &HttpService{
		srv:       &srv,
		tls:       tls,
		certFile:  certFile,
		keyFile:   keyFile,
		generator: generator,
	}
}

func (s *HttpService) Start() {
	if s.tls {
		golog.Info("Listening for HTTPS requests on", s.srv.Addr)
		s.srv.ListenAndServeTLS(s.certFile, s.keyFile)
	} else {
		golog.Info("Listening for HTTP requests on", s.srv.Addr)
		s.srv.ListenAndServe()
	}
}

func (s *HttpService) Stop() {
	if s.srv == nil {
		return
	}

	if err := s.srv.Shutdown(context.Background()); err != nil {
		golog.Errorf("An error occured while shutting down http server: %s", err.Error())
		s.srv.Close()
	}

	s.srv = nil
}
