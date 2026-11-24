package main

import (
	"github.com/Ayushmangit/mirrormatev1/internal/env"
	"go.uber.org/zap"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync() // important to flush the logger for buffered entries

	app := application{
		config: cfg,
		logger: logger,
	}
	if err := app.run(app.mount()); err != nil {
		app.logger.Infow("Error starting the server", "error", err)
	}
}
