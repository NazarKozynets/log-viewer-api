package main

import (
	"log"

	"github.com/nazarkozynets/log-viewer-api/internal/auth"
	"github.com/nazarkozynets/log-viewer-api/internal/config"
	"github.com/nazarkozynets/log-viewer-api/internal/logs"
	"github.com/nazarkozynets/log-viewer-api/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logsService := logs.NewService(
		cfg.LogFiles,
		cfg.MaxReadBytes,
		cfg.DefaultLimit,
		cfg.MaxLimit,
	)

	logsHandler := logs.NewHandler(logsService)

	authClient := auth.NewClient(cfg.AuthMeURL, cfg.AuthLoginURL)
	authHandler := auth.NewHandler(authClient)

	router := server.NewRouter(logsHandler, authHandler)
	handler := server.WithCORS(router, cfg.CORSOrigin)
	app := server.New(cfg.Addr(), handler)

	log.Printf("log-viewer-api started on %s", cfg.Addr())

	if err := app.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
