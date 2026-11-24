package main

import (
	"github.com/Ayushmangit/mirrormatev1/internal/db"
	"github.com/Ayushmangit/mirrormatev1/internal/env"
	"github.com/Ayushmangit/mirrormatev1/internal/store"
	"go.uber.org/zap"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			dsn:             env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/mirrormate?sslmode=disable"),
			maxOpenConns:    env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns:    env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxConnIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync() // important to flush the logger for buffered entries

	db, err := db.New(cfg.db.dsn, cfg.db.maxConnIdleTime, cfg.db.maxIdleConns, cfg.db.maxOpenConns)
	defer db.Close()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("connection to the database established")

	store := store.NewStorage(db)
	app := application{
		config: cfg,
		logger: logger,
		store:  store,
	}
	if err := app.run(app.mount()); err != nil {
		app.logger.Infow("Error starting the server", "error", err)
	}
}
