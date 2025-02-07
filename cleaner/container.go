package cleaner

import (
	"context"
	"fmt"

	"github.com/lucasmendesl/beerus/docker"
	"golang.org/x/sync/errgroup"
)

// https://mohit8830.medium.com/mastering-docker-restart-policies-a-comprehensive-guide-ade2260e7e2c

// listAllowedContainersToRemove returns a list of containers that are allowed to be removed.
// A container is removable if it is either dead or exited and has no restart policy,
// or if it is in created status and has been in this status for more than 5 minutes,
// or if it is in on-failure restart policy and the maximum retry count has been reached.
// The function takes a context.Context and returns a slice of strings containing the IDs of the
// removable containers, and an error if any occurs.
func (c *cleaner) listAllowedContainersToRemove(ctx context.Context) ([]string, error) {
	// Fetch the containers that are either dead or exited and have no restart policy.
	// The containers that are in created status are also considered for removal.
	containers, err := c.d.ListContainers(
		ctx,
		c.config.ConcurrencyLevel,
		docker.WithContainerStatus(
			docker.ContainerStatusDead,
			docker.ContainerStatusExited,
			docker.ContainerStatusCreated,
		),
		docker.WithContainerLabel(c.config.Containers.IgnoreLabels...),
	)

	if err != nil {
		return nil, err
	}

	// Filter the containers that are removable.
	removableContainers := make([]string, 0, len(containers))
	for _, ctr := range containers {
		if docker.CanRemoveContainer(ctr, c.config.Containers.MaxAlwaysRestartPolicyCount) {
			removableContainers = append(removableContainers, ctr.ID)
		}
	}

	return removableContainers, nil
}

// removeContainers removes the specified Docker containers concurrently.
// It logs the start of the removal process and attempts to remove each
// container by calling the Docker API.
// If an error occurs during the removal of any container, it continues the
// operation and logs the error. The function blocks until all containers
// have been processed or the context is canceled.
//
// Parameters:
// - ctx: The context for managing request lifetime and cancellation.
// - containers: A slice of strings containing the IDs of the containers to be removed.
func (c *cleaner) removeContainers(ctx context.Context, containers ...string) error {
	containersLen := len(containers)

	if containersLen == 0 {
		c.log.Warn("No containers to remove")
		return nil
	}

	c.log.Debug("Starting to remove containers...", "count", containersLen)
	g, ctx := errgroup.WithContext(ctx)

	for _, container := range containers {
		g.Go(func() error {
			c.log.Debug("Attempting to remove container", "containerID", container)
			removeOptions := docker.RemoveContainerOptions{
				ContainerID:   container,
				RemoveVolumes: c.config.Containers.ForceVolumeCleanup,
				RemoveLinks:   c.config.Containers.ForceLinkCleanup,
			}

			if err := c.d.RemoveContainer(ctx, removeOptions); err != nil {
				return fmt.Errorf("error removing container with id %s: %w", container, err)
			}

			c.log.Debug("Successfully removed container", "containerID", container)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	c.log.Debug("Successfully removed all containers")
	return nil
}
