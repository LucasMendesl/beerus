//go:generate go run go.uber.org/mock/mockgen -source client.go -package docker -destination ./mocks/client.go Client

package docker

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types/image"
)

type Client interface {
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)

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
