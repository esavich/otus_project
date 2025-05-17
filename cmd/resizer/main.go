package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/esavich/otus_project/internal/config"
	"github.com/esavich/otus_project/internal/logger"
	"github.com/esavich/otus_project/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	logger.SetupLogger(cfg)
	slog.Debug(fmt.Sprintf("config: %+v", cfg))
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	srv := server.NewServer(cfg)

	go func() {
		slog.Debug("Starting server")
		err := srv.Start()
		if err != nil {
			slog.Error(fmt.Sprintf("Error starting server: %s", err))
			cancel()
			return
		}
	}()

	<-ctx.Done()
	slog.Info("Received shutdown signal, shutting down...")
	cancel()
}
