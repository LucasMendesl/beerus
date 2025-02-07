package cleaner

import (
	"context"
	"log/slog"

	"github.com/lucasmendesl/beerus/config"
	"github.com/lucasmendesl/beerus/docker"
)

type cleaner struct {
	d      docker.BeerusContainerAPI
	config *config.Beerus
	log    *slog.Logger
}

// New returns a new cleaner object that can be used to remove images and
// containers that are marked for removal and set up event watchers for
// image untag and container exit events. The function takes a docker
// client and a configuration object as parameters and returns the
// cleaner object.
func New(
	d docker.BeerusContainerAPI,
	config *config.Beerus,
	log *slog.Logger,
) *cleaner {
	return &cleaner{
		d:      d,
		config: config,
		log:    log,
	}
}

// Run starts the cleaner, which removes images and containers that are
// marked for removal and sets up event watchers for image untag and
// container exit events. The function takes a context.Context and
// returns an error if any occurs during the cleanup process. The
// function will block until the context is canceled and will return the
// context's error in this case.
func (c *cleaner) Run(ctx context.Context) error {
	c.log.Info("Starting cleaner, listing containers allowed for removal")
	containers, err := c.listAllowedContainersToRemove(ctx)
	if err != nil {
		c.log.Error("Failed to list removable containers", "error", err)
		return err
	}

	c.log.Info("Removing containers", "count", len(containers))
	if err := c.removeContainers(ctx, containers...); err != nil {
		c.log.Error("Failed to remove containers", "error", err)
		return err
	}

	c.log.Info("Listing images allowed for removal")
	images, err := c.listAllowedImagesToRemove(ctx)
	if err != nil {
		c.log.Error("Failed to list removable images", "error", err)
		return err
	}

	c.log.Info("Removing images", "count", len(images))
	if err := c.removeImages(ctx, images...); err != nil {
		c.log.Error("Failed to remove images", "error", err)
		return err
	}

	c.log.Info("Setting up event watchers")
	workerErr := make(chan error, 1)
	defer func() {
		close(workerErr)
		c.d.Close()
	}()

	go c.watch(ctx, workerErr)

	for {
		select {
		case <-ctx.Done():
			c.log.Info("process canceled, shutting down cleaner")
			return ctx.Err()
		case err := <-workerErr:
			c.log.Error("Error occurred in worker", "error", err)
			return err
		}
	}
}
