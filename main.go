package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Atharva21/streakr/cmd"
	_ "github.com/Atharva21/streakr/internal/streakr"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		<-sigChan
		cancel() // Cancel the context on interrupt
	}()
	slog.Info("Starting Streakr application...")

	cmd.Execute(ctx)
}
