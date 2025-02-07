package cleaner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucasmendesl/beerus/docker"
	"golang.org/x/sync/errgroup"
)

// listAllowedImagesToRemove returns a list of Docker images that are considered
// removable based on specific criteria. It fetches all images and filters
// them to identify those that are either dangling or expired according to
// the provided lifetime threshold. It then removes images that are currently
// running from the list of removable images. The function takes a context.Context
// and returns a slice of docker.Image containing removable images and an error
// if any occurs during the cleanup process.
func (c *cleaner) listAllowedImagesToRemove(ctx context.Context) ([]docker.Image, error) {
	c.log.Debug("Listing allowed images for removal")
	containers, err := c.d.ListContainers(ctx,
		c.config.ConcurrencyLevel,
		docker.WithContainerStatus(docker.ContainerStatusRunning),
	)

	if err != nil {
		return nil, fmt.Errorf("error listing running containers: %w", err)
	}

	c.log.Debug("Removing running images, even if they are expired")
	runningImages := make(map[string]struct{})
	for _, container := range containers {
		runningImages[container.ImageID] = struct{}{}
	}

	c.log.Debug("Getting expired images")
	expiredImgs, err := c.d.ListExpiredImages(ctx, docker.ExpiredImageListOptions{
		LifetimeThresholdInDays: c.config.Images.LifetimeThreshold,
		IgnoreLabels:            c.config.Images.IgnoreLabels,
	})

	if err != nil {
		return nil, fmt.Errorf("error getting expired images: %w", err)
	}

	c.log.Debug("Filtering running images from expired images")
	removableImgs := make([]docker.Image, 0, len(expiredImgs))
	for _, img := range expiredImgs {
		if _, ok := runningImages[img.ID]; !ok {
			removableImgs = append(removableImgs, img)
		}
	}

	c.log.Debug("Returning removable images", "count", len(removableImgs))
	return removableImgs, nil
}

// removeImages removes the specified Docker images concurrently.
// It logs the start of the removal process and attempts to remove each
// image by calling the Docker API.
// If an error occurs during the removal of any image, it continues the
// operation and logs the error. The function blocks until all images
// have been processed or the context is canceled.
//
// Parameters:
// - ctx: The context for managing request lifetime and cancellation.
// - removableImgs: A slice of docker.Image containing the images to be removed.
func (c *cleaner) removeImages(ctx context.Context, removableImgs ...docker.Image) error {
	imagesLen := len(removableImgs)
	c.log.Debug("Removing images", "count", imagesLen)

	if imagesLen == 0 {
		slog.Warn("No images to remove")
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, img := range removableImgs {
		g.Go(func() error {
			c.log.Debug("Attempting to remove image", "imageID", img.ID)
			if err := c.d.RemoveImage(ctx, img.ID); err != nil {
				return fmt.Errorf("error removing image with id %s: %w", img.ID, err)
			}
			c.log.Debug("Successfully removed image", "imageID", img.ID)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
