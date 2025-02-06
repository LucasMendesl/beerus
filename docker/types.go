//go:generate go run go.uber.org/mock/mockgen -source types.go -package docker -destination ./mocks/beerus_container_api.go BeerusContainerAPI

package docker

import (
	"context"
)

type BeerusContainerAPI interface {
	ListExpiredImages(ctx context.Context, options ExpiredImageListOptions) ([]Image, error)
	RemoveImage(ctx context.Context, dockerImage string) error
	Close() error
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
