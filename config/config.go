package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Logging struct {
	// Level defines the level of logging verbosity for the application.
	// Possible values could include "debug", "info", "warn", "error", or "fatal".
	// It controls the amount of logging output generated during the application's execution.
	Level string `mapstructure:"level"`

	// Format defines the format of the log output for the application.
	// Possible values could include "json", "text", or other custom formats.
	// It determines how log messages are structured in the application's logs.
	Format string `mapstructure:"format"`
}

type Beerus struct {
	// Logging specifies the logging configuration, including log level and format.
	Logging Logging `mapstructure:"logging"`
}

// Config represents configuration settings for managing Docker images and containers.
// It contains version information and nested structures for various configuration categories.
type Config struct {
    // Version specifies the version of the application.
    // This field represents the current version of the application
    // It is used to track changes, updates, and compatibility of
    // the application over time.
    Version string `mapstructure:"version"`

	// Beerus holds the configuration settings specific to the Beerus application.
	// It includes settings, logging, images, and container-related configurations.
	Beerus *Beerus `mapstructure:"beerus"`
}

// Load returns a pointer to a Config struct with default values.
// It is used to load default configuration settings for the application.
func Load(configFile string) *Config {
	var config Config

	viper.SetConfigType("yaml")
	viper.SetConfigName("beerus")

	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/beerus")

	viper.SetEnvPrefix("beerus")
	viper.AutomaticEnv()

	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("Failed to parse configuration file")
	} else {
		slog.Info("Using configuration file", "file", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("Failed to unmarshal configuration", "error", err)
	}

	return &config
}
