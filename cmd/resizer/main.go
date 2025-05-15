package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/esavich/otus_project/internal/config"
	"github.com/esavich/otus_project/internal/server"
)

func main() {
	cfg := config.Load()
	fmt.Println("Config loaded:", cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	fmt.Println("starting server with ctx", ctx)
	srv := server.NewServer(cfg)

	go func() {
		err := srv.Start()
		if err != nil {
			fmt.Println("Error starting server:", err)
			cancel()
			return
		}
	}()

	<-ctx.Done()
	fmt.Println("Received shutdown signal, shutting down...")
	cancel()
}
