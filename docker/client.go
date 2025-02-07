//go:generate go run go.uber.org/mock/mockgen -source client.go -package docker -destination ./mocks/client.go Client

package docker

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/image"
)

type Client interface {
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	Events(ctx context.Context, options events.ListOptions) (<-chan events.Message, <-chan error)

	Ping(ctx context.Context) (types.Ping, error)
	Close() error
}

type dockerClient struct {
	cli Client
	log *slog.Logger
}

// New returns a new Client instance that can be used to interact with the Docker
// engine.
func New(cli Client, logger *slog.Logger) BeerusContainerAPI {
	return &dockerClient{cli: cli, log: logger}
}

// Close closes the Docker client connection, releasing any resources that
// were allocated. It returns an error if there is an issue during the
// closure process.
func (d *dockerClient) Close() error {
	return d.cli.Close()
}
