package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Atharva21/streakr/internal/shutdown"
)

var bootstrapConfigOnce sync.Once

type StreakrConfig struct {
	ConfigRootDir string
	DataDir       string
	LogFileDir    string
	LogFileName   string
	StoreName     string
}

var streakrConfigInstance *StreakrConfig = nil

func GetStreakrConfig() StreakrConfig {
	if streakrConfigInstance == nil {
		panic("StreakrConfig is not initialized. Call BootstrapConfig() first.")
	}
	return *streakrConfigInstance
}

func exitWithStderrGeneric(err error) {
	fmt.Fprintf(os.Stderr, "An unexpected error occurred while bootstrapping necessary config files for streakr: %s", err.Error())
	shutdown.GracefulShutdown(1)
}

func BootstrapConfig() {
	bootstrapConfigOnce.Do(func() {
		streakrConfigInstance = &StreakrConfig{}
		userHomeDir, err := os.UserConfigDir()
		if err != nil {
			exitWithStderrGeneric(err)
		}
		streakrConfigInstance.ConfigRootDir = filepath.Join(userHomeDir, "streakr")
		streakrConfigInstance.DataDir = filepath.Join(streakrConfigInstance.ConfigRootDir, "data")
		streakrConfigInstance.LogFileDir = filepath.Join(streakrConfigInstance.ConfigRootDir, "logs")
		streakrConfigInstance.LogFileName = "streakr.log"
		streakrConfigInstance.StoreName = "streakr.db"

		// Create necessary directories
		err = os.MkdirAll(streakrConfigInstance.ConfigRootDir, 0700)
		if err != nil {
			exitWithStderrGeneric(err)
		}
		err = os.MkdirAll(streakrConfigInstance.DataDir, 0700)
		if err != nil {
			exitWithStderrGeneric(err)
		}
		err = os.MkdirAll(streakrConfigInstance.LogFileDir, 0700)
		if err != nil {
			exitWithStderrGeneric(err)
		}
	})
}
