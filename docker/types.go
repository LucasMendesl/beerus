//go:generate go run go.uber.org/mock/mockgen -source types.go -package docker -destination ./mocks/beerus_container_api.go BeerusContainerAPI

package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
)

type BeerusContainerAPI interface {
	Inspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ListContainers(ctx context.Context, concurrency uint8, options ...ListContainersOptions) ([]Container, error)
	RemoveContainer(ctx context.Context, options RemoveContainerOptions) error
	ListExpiredImages(ctx context.Context, options ExpiredImageListOptions) ([]Image, error)
	RemoveImage(ctx context.Context, dockerImage string) error
	FromEvents(ctx context.Context, actions ...events.Action) <-chan EventResult
	Close() error
}

type ListContainersOptions func(*ListContainersParams)

// ContainerStatus is a string type representing the status of a container. It
// has one of the following values:
//
//   - running: the container is currently running.
//   - exited: the container has exited, but its process has not been removed from
//     memory yet.
//   - dead: the container was only partially removed because resources were kept
//     busy by an external process.
type ContainerStatus string

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
)

// ExpiredImageListOptions represents criteria for removable images.
type ExpiredImageListOptions struct {
	LifetimeThresholdInDays uint16
	IgnoreLabels            []string
}

// RemoveContainerOptions represents options for removing a container.
type RemoveContainerOptions struct {
	ContainerID   string
	RemoveVolumes bool
	RemoveLinks   bool
}

// Container represents a Docker container, containing its ID, status, image name,
// and image ID.
type Container struct {
	ID            string
	Image         string
	ImageID       string
	Labels        map[string]string
	CreatedAt     time.Time
	Status        ContainerStatus
	RestartCount  int
	RestartPolicy container.RestartPolicy
}

// Image represents a Docker image, containing its ID, tags, and labels.
type Image struct {
	ID     string
	Labels map[string]string
	Tags   []string
}

// EventResult represents a result from the event stream, which may contain
// either a Message or an error.
type EventResult struct {
	Message events.Message
	Err     error
}

type ListContainersParams struct {
	Status []ContainerStatus
	Label  []string
}
