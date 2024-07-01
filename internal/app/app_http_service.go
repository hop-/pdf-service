package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hop-/golog"
)

func (a *App) startHttp() {
	defer a.wg.Done()

	addr := fmt.Sprintf(":%d", a.options.Http.Port)

	var router http.Handler
	// TODO: add router

	srv := http.Server{Addr: addr, Handler: router}
	a.OnShutdown(func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			golog.Errorf("An error occured while shutting down http server: %s", err.Error())
			srv.Close()
		}
	})

	if a.options.Http.Secure {
		golog.Info("Listening for HTTPS requests on", addr)
		srv.ListenAndServeTLS(a.options.Http.Cert, a.options.Http.Key)
	} else {
		golog.Info("Listening for HTTP requests on", addr)
		srv.ListenAndServe()
	}
}
