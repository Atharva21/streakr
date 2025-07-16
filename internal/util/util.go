package util

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/Atharva21/streakr/internal/shutdown"
)

var (
	logFileAbsolutePath string
	bootstrapUtilOnce   sync.Once
)

func BootstrapUtil(logLocation string) {
	bootstrapUtilOnce.Do(func() {
		logFileAbsolutePath = logLocation
	})
}

func ErrorAndExitGeneric(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
	fmt.Fprintf(os.Stderr, "An unexpected error occurred. Please check the logs at %s for more details.\n", logFileAbsolutePath)
	shutdown.GracefulShutdown(1)
}

func ErrorAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	shutdown.GracefulShutdown(1)
}
