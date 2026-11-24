package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type application struct {
	config config
	logger *zap.SugaredLogger
}

type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}
	app.logger.Infow("Server is listening", "port", app.config.addr)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
