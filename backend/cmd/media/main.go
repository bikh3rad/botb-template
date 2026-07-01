package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// main is the entrypoint for the media service binary.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, err := wireApp(ctx)
	if err != nil {
		panic(err)
	}

	logger := application.GetLogger().With("component", "main", "service", "media")
	slog.SetDefault(logger)

	logger.Info("media service starting...")

	if err := application.Start(ctx); err != nil {
		panic(err)
	}
	defer application.Shutdown(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("media service stopping...", "signal", sig)
}
