package util

import (
	"fmt"
	"os"

	"github.com/Atharva21/streakr/internal/config"
	"github.com/Atharva21/streakr/internal/shutdown"
)

func ErrorAndExitGeneric() {
	fmt.Fprintf(os.Stderr, "An unexpected error occurred. Please check the logs at %s for more details.\n", config.GetStreakrConfig().LogFileDir)
	shutdown.GracefulShutdown(1)
}

func ErrorAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	shutdown.GracefulShutdown(1)
}
