package taipan

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"sync"
)

var globalConfig *viper.Viper
var globalConfigMutex = &sync.Mutex{}

// Get returns the configured global config object
// It's created the first time Get is called
// For example profiles="profile1,profile2" resolves files in the following order:
// config.yaml, config-profile1.yaml, config-profile2.yaml, config-local.yaml
// Panics if file not found, except for missing config-local.yaml which is ignored
func Get() *viper.Viper {
	// Make sure globalConfig is only initialized once
	globalConfigMutex.Lock()
	defer globalConfigMutex.Unlock()

	// Return if already initialized
	if globalConfig != nil {
		return globalConfig
	}

	globalConfig = viper.New()

	// Read configurable config path
	configPath := os.Getenv("TAIPAN_CONFIG_PATH")
	if configPath == "" {
		configPath = "resources"
	}

	globalConfig.AddConfigPath(configPath)

	// Set defaults
	globalConfig.SetConfigName("config")
	globalConfig.SetConfigType("yaml")
	globalConfig.SetEnvPrefix("TPN")
	globalConfig.SetEnvKeyReplacer(strings.NewReplacer(
		".", "_",
		"-", "_"))
	globalConfig.AutomaticEnv()

	// Read <configPath>/config.yaml, panicking on failure
	err := globalConfig.ReadInConfig()
	if err != nil {
		panic("Failed to load base config file")
	}

	// Get comma-separated profiles. Merge them in in order overriding previous values
	// A profile resolves to a file like <configPath>/config-<profile>.yml
	profiles := os.Getenv("TAIPAN_PROFILES")
	if profiles != "" {
		splitProfiles := strings.Split(profiles, ",")
		for _, profile := range splitProfiles {
			profileLowerTrimmed := strings.ToLower(strings.TrimSpace(profile))
			globalConfig.SetConfigName(fmt.Sprintf("config-%s", profileLowerTrimmed))
			err := globalConfig.MergeInConfig()
			if err != nil {
				panic(fmt.Sprintf("failed to load config file for profile %s: %s\n", profileLowerTrimmed, err))
			}
		}
	}

	// Merge in ./config-local.yml if present
	globalConfig.AddConfigPath(".")
	globalConfig.SetConfigName("config-local")
	err = globalConfig.MergeInConfig()
	if err == nil {
		log.Printf("Loaded local config")
	}

	return globalConfig
}
