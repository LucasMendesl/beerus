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

type Image struct {
	// LifetimeThreshold represents the threshold in terms of time (in days)
	// after which images are considered for removal. Images older than this threshold may be cleaned up.
	LifetimeThreshold uint16 `mapstructure:"lifetimeThreshold"`

	// IgnoreLabels contains a list of image labels that should be ignored during the cleanup process.
	// Images with any of these labels will not be considered for removal.
	IgnoreLabels []string `mapstructure:"ignoreLabels"`

	// ForceRemovalOnConflict is a boolean that, if set to true, will force the
	// removal of resources when a conflict is detected during the cleanup
	// process (when one repository have more than one tag).
	ForceRemovalOnConflict bool `mapstructure:"forceRemovalOnConflict"`
}

type Container struct {
	// MaxAlwaysRestartPolicyCount defines the maximum number of times a container can be restarted
	// within a specific time window before it is considered for removal.
	// If a container restarts more than this number of times, it is considered to be in a
	// restart loop and will be removed to prevent resource waste (using restart policy always).
	MaxAlwaysRestartPolicyCount int `mapstructure:"maxAlwaysRestartPolicyCount"`

	// IgnoreLabels contains a list of image labels that should be ignored during the cleanup process.
	// Images with any of these labels will not be considered for removal.
	IgnoreLabels []string `mapstructure:"ignoreLabels"`

	// ForceVolumeCleanup is a boolean that, if set to true, will force the removal of volumes associated with
	// containers that are being removed. This can be useful for cleaning up volumes that are no longer in
	// use, but it may also cause loss of data if volumes are being used by other containers.
	// By default, volumes are not removed when a container is removed, in order to prevent data loss.
	// However, if a container is being removed due to a restart loop, and the container is configured to
	// always restart, then the volume will be removed to prevent resource waste.
	ForceVolumeCleanup bool `mapstructure:"forceVolumeCleanup"`

	// ForceLinkCleanup is a boolean that, if set to true, will force the removal of links
	// associated with containers that are being removed. This can be useful for cleaning up
	// links that are no longer in use, but it may also cause loss of connectivity if links
	// are being used by other containers. By default, links are not removed when a container
	// is removed, to prevent connectivity issues. However, if a container is being removed
	// due to a restart loop, and it is configured to always restart, then the link will be
	// removed to prevent resource waste.
	ForceLinkCleanup bool `mapstructure:"forceLinkCleanup"`
}

type Beerus struct {
	// ConcurrencyLevel defines the maximum number of goroutines that can run in parallel
	// during the execution of the application. It controls how many items
	// are processed at the same time. A higher value can lead to faster cleaning but
	// may also increase the load on the system.
	ConcurrencyLevel uint8 `mapstructure:"concurrencyLevel"`

	// ExpirePollCheckInterval specifies the interval in hours between each poll
	// check for expired images. It controls how frequently the application will
	// check for images that are older than the ImageLifetimeThreshold value.
	// A higher value can lead to less frequent checks and lower system load,
	// but may also mean expired images are removed less quickly.
	ExpirePollCheckInterval uint8 `mapstructure:"expiringPollCheckInterval"`

	// Logging specifies the logging configuration, including log level and format.
	Logging Logging `mapstructure:"logging"`

	// Images contains settings related to Docker image management, such as lifetime thresholds.
	Images Image `mapstructure:"images"`

	// Containers includes configuration parameters for managing Docker containers,
	// particularly related to restart policies and removal criteria.
	Containers Container `mapstructure:"containers"`
}

// Config represents configuration settings for managing Docker images and containers.
// It contains version information and nested structures for various configuration categories.
type Config struct {
	// Version specifies the version of the configuration file format.
	// It is used to handle changes to the configuration file format
	// over time and to ensure backwards compatibility.
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
