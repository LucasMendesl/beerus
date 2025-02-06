package docker

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/docker/docker/api/types/image"
)

const (
	danglingImageTag = "<none>:<none>"
	dayInHours       = 24
)

// ListExpiredImages retrieves a list of Docker images that are considered
// removable based on specific criteria. It fetches all images and filters
// them to identify those that are either dangling or expired according to
// the provided lifetime threshold.
//
// Parameters:
//   - ctx: The context for managing request lifetime and cancellation.
//   - options: A struct containing criteria for removable images, specifically
//     the lifetime threshold in days.
//
// Returns:
//   - A slice of image.Summary containing removable images.
//   - An error if there is an issue retrieving the list of images.
func (d *dockerClient) ListExpiredImages(ctx context.Context, options ExpiredImageListOptions) ([]Image, error) {
	images, err := d.cli.ImageList(ctx, image.ListOptions{
		All: true,
	})

	if err != nil {
		return nil, fmt.Errorf("expired docker images error: %w", err)
	}

	removableImages := make([]Image, 0, len(images))

	if len(images) == 0 {
		return removableImages, nil
	}

	for _, image := range images {
		isDangling := slices.Contains(image.RepoTags, danglingImageTag)
		imageExpired := isImageExpired(image.Created, options.LifetimeThresholdInDays)

		if isDangling || imageExpired {
			removableImages = append(removableImages, Image{
				ID:     image.ID,
				Labels: image.Labels,
				Tags:   image.RepoTags,
			})
		}
	}

	return removeIgnored(removableImages, options.IgnoreLabels...), nil
}

// RemoveImage removes a Docker image by its ID or name.
//
// Parameters:
//   - ctx: The context for managing request lifetime and cancellation.
//   - dockerImage: The ID or name of the image to be removed.
//
// Returns:
//   - An error if there is an issue removing the image.
func (d *dockerClient) RemoveImage(ctx context.Context, dockerImage string) error {
	_, err := d.cli.ImageRemove(ctx, dockerImage, image.RemoveOptions{})
	return err
}

// isImageExpired checks if a Docker image is expired based on its creation
// time and the given lifetime threshold in days.
//
// Parameters:
//   - created: The creation time of the image in seconds since the Unix
//     epoch.
//   - lifetimeThresholdInDays: The lifetime threshold in days.
//
// Returns:
//   - A boolean indicating whether the image is expired or not.
func isImageExpired(created int64, lifetimeThresholdInDays uint16) bool {
	createdTime := time.Unix(created, 0)
	days := uint16(time.Since(createdTime).Hours() / dayInHours)

	return days >= lifetimeThresholdInDays
}
