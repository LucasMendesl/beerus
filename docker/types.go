//go:generate go run go.uber.org/mock/mockgen -source types.go -package docker -destination ./mocks/beerus_container_api.go BeerusContainerAPI

package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/events"
)

type ListContainersOptions func(*ListContainerParams)

// ContainerStatus is a string type representing the status of a container. It
// has one of the following values:
//
//   - running: the container is currently running.
//   - exited: the container has exited, but its process has not been removed from
//     memory yet.
//   - dead: the container was only partially removed because resources were kept
//     busy by an external process.
type ContainerStatus string

type BeerusContainerAPI interface {
	ListContainers(ctx context.Context, concurrency uint8, options ...ListContainersOptions) ([]Container, error)
	ListExpiredImages(ctx context.Context, options ExpiredImageListOptions) ([]Image, error)
	RemoveImage(ctx context.Context, dockerImage string) error
    FromEvents(ctx context.Context, actions ...events.Action) <-chan EventResult
	Close() error
}

// ListContainerParams represents the parameters used when listing Docker containers.
// It includes criteria such as the desired container statuses and labels to filter by.
type ListContainerParams struct {
	Status []ContainerStatus
	Label  []string
}

// ExpiredImageListOptions represents criteria for removable images.
type ExpiredImageListOptions struct {
	LifetimeThresholdInDays uint16
	IgnoreLabels            []string
}

// Image represents a Docker image, containing its ID, tags, and labels.
type Image struct {
	ID     string
	Labels map[string]string
	Tags   []string
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
}

// EventResult represents a result from the event stream, which may contain
// either a Message or an error.
type EventResult struct {
    Message events.Message
    Err     error
}
