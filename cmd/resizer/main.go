package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/esavich/otus_project/internal/config"
	"github.com/esavich/otus_project/internal/diskcache"
	"github.com/esavich/otus_project/internal/logger"
	"github.com/esavich/otus_project/internal/server"
	"github.com/esavich/otus_project/internal/service"
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

	// main dependencies
	// u can use
	imageService := service.NewSimpleImageService()
	dc, err := diskcache.NewDiskCacheWrapper(cfg.Cache.MaxItems, cfg.Cache.Path)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating disk cache: %s", err))
		return
	}
	cachedService := service.NewCachedImageService(imageService, dc)

	srv := server.NewServer(cfg, cachedService)

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
	err = dc.ClearDiskCache()
	if err != nil {
		slog.Error(fmt.Sprintf("Error clearing disk cache: %s", err))
	}
	slog.Info("Cache cleared")

	cancel()
}
