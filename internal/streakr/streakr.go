package streakr

import (
	"path/filepath"
	"sync"

	"github.com/Atharva21/streakr/internal/config"
	"github.com/Atharva21/streakr/internal/log"
)

var bootstrapOnce sync.Once

func bootsrapStreakr() {
	bootstrapOnce.Do(func() {
		// bootstrap app config
		config.BootstrapConfig()
		appConfig := config.GetStreakrConfig()

		// bootstrap logger
		log.BootsrapLogger(filepath.Join(appConfig.LogFileDir, appConfig.LogFileName))
	})
}

func init() {
	bootsrapStreakr()
}
