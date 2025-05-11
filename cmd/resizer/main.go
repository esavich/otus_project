package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

func main() {
	//config.Load()
	//cache.Setup(config.Cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Println("starting server with ctx", ctx)
	//go server.Run()

	<-ctx.Done()
	fmt.Println("Received shutdown signal, shutting down...")
	cancel()

}
