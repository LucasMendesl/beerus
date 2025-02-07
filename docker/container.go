package docker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

const statusFilter = "status"

// WithContainerStatus filters containers by status when calling ListContainers.
// It returns a ListContainersOptions that sets the Status field of the
// ListContainersParams struct. It takes a variable number of arguments of type
// ContainerStatus, and multiple statuses can be passed to filter containers by
// multiple statuses.
func WithContainerStatus(status ...ContainerStatus) ListContainersOptions {
	return func(o *ListContainersParams) {
		o.Status = status
	}
}

// WithContainerLabel filters containers by label when calling ListContainers.
// It returns a ListContainersOptions that sets the Label field of the
// ListContainersParams struct. It takes a variable number of arguments of type
// string, allowing multiple labels to be passed to filter containers by
// multiple labels.
func WithContainerLabel(label ...string) ListContainersOptions {
	return func(o *ListContainersParams) {
		o.Label = label
	}
}

// ListContainers retrieves a list of Docker containers based on their status.
// It creates concurrent goroutines to inspect each container and fetch its
// details. The function filters containers by status and returns a slice of
// Container objects containing their IDs, images, labels, creation time, and
// current status.
//
// Parameters:
// - ctx: The context for managing request lifetime and cancellation.
// - concurrency: The number of concurrent goroutines to use for inspecting containers.
// - status: Variadic parameter specifying one or more ContainerStatus values to filter the containers.
//
// Returns:
// - A slice of Container objects with details of each container.
// - An error if there is an issue fetching or inspecting the containers.
func (d *dockerClient) ListContainers(ctx context.Context, concurrency uint8, options ...ListContainersOptions) ([]Container, error) {
	listContainerParam := &ListContainersParams{
		Status: []ContainerStatus{},
		Label:  []string{},
	}
	for _, option := range options {
		option(listContainerParam)
	}

	containerFilters := filters.NewArgs()
	for _, s := range listContainerParam.Status {
		containerFilters.Add(statusFilter, string(s))
	}

	listOptionsParams := container.ListOptions{
		All:     true,
		Filters: containerFilters,
	}

	containers, err := d.cli.ContainerList(ctx, listOptionsParams)
	if err != nil {
		return nil, fmt.Errorf("fetching containers error: %w", err)
	}

	containersLen := len(containers)
	if containersLen == 0 {
		return []Container{}, nil
	}

	var wg sync.WaitGroup

	containerCh := make(chan Container, concurrency)
	containerList := make([]Container, 0, containersLen)
	filteredContainers := make([]Container, 0)

	for _, c := range containers {
		filteredContainers = append(filteredContainers, Container{
			ID:        c.ID,
			Image:     c.Image,
			ImageID:   c.ImageID,
			Labels:    c.Labels,
			CreatedAt: time.Unix(c.Created, 0),
			Status:    ContainerStatus(c.Status),
		})
	}

	filteredContainers = removeIgnored(filteredContainers, listContainerParam.Label...)
	for _, c := range filteredContainers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			details, err := d.Inspect(ctx, c.ID)
			if err != nil {
				d.log.Error("Failed to inspect container", "error", err, "id", c.ID)
				return
			}

			containerCh <- Container{
				ID:            c.ID,
				Image:         c.Image,
				ImageID:       c.ImageID,
				Labels:        c.Labels,
				CreatedAt:     c.CreatedAt,
				Status:        c.Status,
				RestartCount:  details.RestartCount,
				RestartPolicy: details.HostConfig.RestartPolicy,
			}
		}()
	}

	// putting waitgroup to wait in a separate goroutine
	// so we don't await all events to be processed to read the channel
	go func() {
		wg.Wait()
		close(containerCh)
	}()

	for container := range containerCh {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			containerList = append(containerList, container)
		}
	}

	return containerList, nil
}

// Inspect retrieves detailed information about a Docker container by its ID.
//
// Parameters:
//   - ctx: The context for managing request lifetime and cancellation.
//   - containerID: The ID of the container to be inspected.
//
// Returns:
//   - A types.ContainerJSON object containing detailed information about the container.
//   - An error if there is an issue retrieving the container information.
func (d *dockerClient) Inspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return d.cli.ContainerInspect(ctx, containerID)
}

// RemoveContainer removes a Docker container by its ID.
//
// Parameters:
//   - ctx: The context for managing request lifetime and cancellation.
//   - options: A RemoveContainerOptions object containing the ID of the container
//     to be removed, and flags to indicate if the container's volumes
//     and links should also be removed.
//
// Returns:
//   - An error if there is an issue removing the container.
func (d *dockerClient) RemoveContainer(ctx context.Context, options RemoveContainerOptions) error {
	return d.cli.ContainerRemove(ctx, options.ContainerID, container.RemoveOptions{
		RemoveVolumes: options.RemoveVolumes,
		RemoveLinks:   options.RemoveLinks,
	})
}

// CanRemoveContainer determines if a given Docker container is eligible for removal
// based on its restart policy and count. The function evaluates the container's
// restart policy and checks if the restart count exceeds the specified thresholds.
//
// Parameters:
//   - c: A Container object representing the Docker container to be evaluated.
//   - alwaysRestartThreshold: An integer threshold for containers with an 'always'
//     restart policy to determine if they can be removed.
//
// Returns:
//   - A boolean indicating whether the container can be removed (true) or not (false).
//     It returns true if the container's restart policy is either disabled or unless-stopped,
//     or if the container's restart count meets or exceeds the specified thresholds for
//     'always' or 'on-failure' policies.
func CanRemoveContainer(c Container, alwaysResartThreshold int) bool {
	switch c.RestartPolicy.Name {
	case container.RestartPolicyDisabled, container.RestartPolicyUnlessStopped:
		return true
	case container.RestartPolicyAlways:
		return c.RestartCount >= alwaysResartThreshold && alwaysResartThreshold > 0
	case container.RestartPolicyOnFailure:
		return c.RestartCount >= c.RestartPolicy.MaximumRetryCount
	default:
		return false
	}
}
