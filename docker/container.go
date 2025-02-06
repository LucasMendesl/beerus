package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

const (
	// ContainerStatusRunning means that the container is currently running.
	ContainerStatusRunning ContainerStatus = "running"

	// ContainerStatusExited means that the container has exited, but its
	// process has not been removed from memory.
	ContainerStatusExited ContainerStatus = "exited"

	// ContainerStatusDead menans that the container that was only partially
	// removed because resources were kept busy by an external process.
	ContainerStatusDead ContainerStatus = "dead"

	// ContainerStatusCreated means that the container has been created but not yet started.
	ContainerStatusCreated ContainerStatus = "created"

    statusFilter = "status"
)

// WithContainerStatus filters containers by status when calling ListContainers.
// It returns a ListContainersOptions that sets the Status field of the
// ListContainersParams struct. It takes a variable number of arguments of type
// ContainerStatus, and multiple statuses can be passed to filter containers by
// multiple statuses.
func WithContainerStatus(status ...ContainerStatus) ListContainersOptions {
	return func(o *ListContainerParams) {
		o.Status = status
	}
}

// WithContainerLabel filters containers by label when calling ListContainers.
// It returns a ListContainersOptions that sets the Label field of the
// ListContainersParams struct. It takes a variable number of arguments of type
// string, allowing multiple labels to be passed to filter containers by
// multiple labels.
func WithContainerLabel(label ...string) ListContainersOptions {
	return func(o *ListContainerParams) {
		o.Label = label
	}
}

func (d *dockerClient) ListContainers(ctx context.Context, concurrency uint8, options ...ListContainersOptions) ([]Container, error) {
	listContainerParam := &ListContainerParams{
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

	return removeIgnored(filteredContainers, listContainerParam.Label...), nil
}
