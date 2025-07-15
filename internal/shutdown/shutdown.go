package shutdown

import (
	"log/slog"
	"os"
)

type CleanupFunc func() error

var cleanupHooks []CleanupFunc = make([]CleanupFunc, 0)
var logCleanupHook CleanupFunc

func RegisterCleanupHook(hook CleanupFunc) {
	if hook != nil {
		cleanupHooks = append(cleanupHooks, hook)
	}
}

func RegisterLogCleanupHook(hook CleanupFunc) {
	if hook != nil {
		logCleanupHook = hook
	}
}

func cleanup() {
	for _, hook := range cleanupHooks {
		if hook == nil {
			continue
		}
		if err := hook(); err != nil {
			// log errors if any cleanup hook fails
			slog.Error("Error during cleanup hook",
				slog.String("error", err.Error()),
			)
		}
	}
	if logCleanupHook != nil {
		// ignore possible errors as we may not have a file to log to?
		logCleanupHook()
	}
}

func GracefulShutdown(exitCode int) {
	cleanup()
	os.Exit(exitCode)
}
