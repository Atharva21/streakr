package config

import (
	"os"
	"path/filepath"
	"sync"
)

var bootstrapConfigOnce sync.Once

type StreakrConfig struct {
	ConfigRootDir string
	DataDir       string
	LogFileDir    string
	LogFileName   string
}

var streakrConfigInstance *StreakrConfig = nil

func GetStreakrConfig() StreakrConfig {
	if streakrConfigInstance == nil {
		panic("StreakrConfig is not initialized. Call BootstrapConfig() first.")
	}
	return *streakrConfigInstance
}

func BootstrapConfig() {
	bootstrapConfigOnce.Do(func() {
		streakrConfigInstance = &StreakrConfig{}
		userHomeDir, err := os.UserConfigDir()
		if err != nil {
			panic("Failed to get user config directory: " + err.Error())
		}
		streakrConfigInstance.ConfigRootDir = filepath.Join(userHomeDir, "streakr")
		streakrConfigInstance.DataDir = filepath.Join(streakrConfigInstance.ConfigRootDir, "data")
		streakrConfigInstance.LogFileDir = filepath.Join(streakrConfigInstance.ConfigRootDir, "logs")
		streakrConfigInstance.LogFileName = "streakr.log"

		// Create necessary directories
		err = os.MkdirAll(streakrConfigInstance.ConfigRootDir, 0700)
		if err != nil {
			panic("Failed to create streakr config directory: " + err.Error())
		}
		err = os.MkdirAll(streakrConfigInstance.DataDir, 0700)
		if err != nil {
			panic("Failed to create streakr data directory: " + err.Error())
		}
		err = os.MkdirAll(streakrConfigInstance.LogFileDir, 0700)
		if err != nil {
			panic("Failed to create streakr log directory: " + err.Error())
		}
	})
}
