package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/routes"
)

type HttpService struct {
	srv      *http.Server
	tls      bool
	certFile string
	keyFile  string
}

func NewHttpService(
	port uint16,
	tls bool,
	certFile string,
	keyFile string,
) *HttpService {
	addr := fmt.Sprintf(":%d", port)

	router := routes.NewRouter()
	srv := http.Server{Addr: addr, Handler: router}

	return &HttpService{
		srv:      &srv,
		tls:      tls,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (s *HttpService) Start() {
	if s.tls {
		golog.Info("Listening for HTTPS requests on", s.srv.Addr)

		err := s.srv.ListenAndServeTLS(s.certFile, s.keyFile)
		if err != nil {
			golog.Errorf("Failed to start HTTPS service: %s", err.Error())
		}
	} else {
		golog.Info("Listening for HTTP requests on", s.srv.Addr)
		err := s.srv.ListenAndServe()
		if err != nil {
			golog.Errorf("Failed to start HTTP service: %s", err.Error())
		}
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
