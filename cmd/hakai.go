package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/client"
	"github.com/lucasmendesl/beerus/cleaner"
	"github.com/lucasmendesl/beerus/config"
	"github.com/lucasmendesl/beerus/docker"
	"github.com/lucasmendesl/beerus/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func newHakaiCmd() *cobra.Command {
	hakaiCmd := &cobra.Command{
		Use:   "hakai",
		Short: "Hakai",
		Run:   help,
		RunE:  cleanResources,
	}

	setupCommandFlags(hakaiCmd.Flags())
	return hakaiCmd
}

func help(cmd *cobra.Command, _ []string) {
	cmd.Help()
}

func setupCommandFlags(commandFlags *pflag.FlagSet) {
	// general flags
	commandFlags.Uint8P("concurrency-level", "c", 5, "number of concurrent workers")
	commandFlags.Uint8P("poll-check-interval", "i", 1, "interval to check for expired resources in hours")

	// log section flags
	commandFlags.String("log-level", "debug", "log level (debug, info, warn, error)")
	commandFlags.String("log-format", "text", "log format (json, text)")

	// image section flags
	commandFlags.Uint16("lifetime-threshold", 100, "lifetime threshold in days")
	commandFlags.StringArray("image-ignore-labels", []string{}, "ignore images with the specified label during cleanup")

	// container section flags
	commandFlags.Int("max-always-restart-policy-count", 0, "max always restart policy count (0 is disabled)")
	commandFlags.StringArray("container-ignore-labels", []string{}, "ignore containers with the specified label during cleanup")
	commandFlags.Bool("force-volume-cleanup", false, "force volume cleanup")
	commandFlags.Bool("force-link-cleanup", false, "force link cleanup")

	bindCommandFlags(commandFlags)
}

func bindCommandFlags(commandFlags *pflag.FlagSet) {
	viper.BindPFlag("beerus.concurrencyLevel", commandFlags.Lookup("concurrency-level"))
	viper.BindPFlag("beerus.pollCheckInterval", commandFlags.Lookup("poll-check-interval"))

	viper.BindPFlag("beerus.logging.level", commandFlags.Lookup("log-level"))
	viper.BindPFlag("beerus.logging.format", commandFlags.Lookup("log-format"))

	viper.BindPFlag("beerus.images.lifetimeThreshold", commandFlags.Lookup("lifetime-threshold"))
	viper.BindPFlag("beerus.images.ignoreLabels", commandFlags.Lookup("image-ignore-labels"))

	viper.BindPFlag("beerus.containers.maxAlwaysRestartPolicyCount", commandFlags.Lookup("max-always-restart-policy-count"))
	viper.BindPFlag("beerus.containers.ignoreLabels", commandFlags.Lookup("container-ignore-labels"))
	viper.BindPFlag("beerus.containers.forceVolumeCleanup", commandFlags.Lookup("force-volume-cleanup"))
	viper.BindPFlag("beerus.containers.forceLinkCleanup", commandFlags.Lookup("force-link-cleanup"))
}

func cleanResources(cmd *cobra.Command, _ []string) error {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return fmt.Errorf("error creating docker client api: %w", err)
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	stopSignal := make(chan os.Signal, 1)

	go func() {
		signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT)
		<-stopSignal
		cancel()
	}()

	filePath := cmd.Flag("config-file").Value
	cfg := config.Load(filePath.String())

	logger, err := logger.Create(cfg.Beerus.Logging)
	if err != nil {
		return fmt.Errorf("error creating logger: %w", err)
	}

	cleaner := cleaner.New(docker.New(cli, logger), cfg.Beerus, logger)

	if err := cleaner.Run(ctx); err != nil {
		return fmt.Errorf("error cleaning resources: %w", err)
	}

	return nil
}
