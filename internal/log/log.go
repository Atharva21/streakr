package log

import (
	"log/slog"
	"sync"

	"github.com/Atharva21/streakr/internal/shutdown"
	"gopkg.in/natefinch/lumberjack.v2"
)

var bootstrapLoggerOnce sync.Once

func BootsrapLogger(absoluteLogfilePath string) {
	bootstrapLoggerOnce.Do(func() {
		lumber := &lumberjack.Logger{
			Filename:   absoluteLogfilePath,
			MaxSize:    10,   // MB
			MaxBackups: 3,    // log.1, log.2, etc.
			MaxAge:     28,   // days
			Compress:   true, // compress old logs
		}

		shutdown.RegisterLogCleanupHook(func() error {
			return lumber.Close()
		})
		loggingHandler := slog.NewTextHandler(
			lumber,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		)
		logger := slog.New(loggingHandler)
		slog.SetDefault(logger)
	})
}
