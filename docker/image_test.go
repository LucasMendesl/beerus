package docker_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/lucasmendesl/beerus/docker"
	mock "github.com/lucasmendesl/beerus/docker/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func nopErr(t *testing.T, err error) bool {
	require.NoError(t, err)
	return false
}

type wantErr func(t *testing.T, err error) bool

func TestDockerClient_ListExpiredImages(t *testing.T) {
	var (
		ctrl         = gomock.NewController(t)
		dockerClient = mock.NewMockClient(ctrl)
		logger       = slog.New(slog.NewJSONHandler(io.Discard, nil))

		listImagesError = errors.New("list images error")
	)
	type args struct {
		ctx     context.Context
		options docker.ExpiredImageListOptions
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		expected  []docker.Image
		wantErr   wantErr
	}{
		{
			name: "error on list images",
			args: args{
				ctx: context.Background(),
				options: docker.ExpiredImageListOptions{
					LifetimeThresholdInDays: 100,
				},
			},
			mockSetup: func() {
				dockerClient.
					EXPECT().
					ImageList(
						gomock.Any(),
						gomock.Any(),
					).
					Return(nil, listImagesError).
					Times(1)
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "expired docker images error: list images error")
				return true
			},
		},
		{
			name: "empty image list",
			args: args{
				ctx: context.Background(),
				options: docker.ExpiredImageListOptions{
					LifetimeThresholdInDays: 100,
				},
			},
			mockSetup: func() {
				dockerClient.
					EXPECT().
					ImageList(
						gomock.Any(),
						gomock.Any(),
					).
					Return([]image.Summary{}, nil).
					Times(1)
			},
			wantErr:  nopErr,
			expected: []docker.Image{},
		},
		{
			name: "images expired or dangling",
			args: args{
				ctx: context.Background(),
				options: docker.ExpiredImageListOptions{
					LifetimeThresholdInDays: 100,
				},
			},
			mockSetup: func() {
				dockerClient.
					EXPECT().
					ImageList(
						gomock.Any(),
						gomock.Any(),
					).
					Return([]image.Summary{
						{
							ID:       "d55c68fb3405",
							Created:  time.Now().Add(-time.Hour * 24 * 10).Unix(),
							RepoTags: []string{"golang:latest"},
							Labels:   map[string]string{},
						},
						{
							ID:       "a76d6a1f0270",
							Created:  time.Now().Add(-time.Hour * 24 * 110).Unix(),
							RepoTags: []string{"nginx:latest"},
							Labels:   map[string]string{},
						},
						{
							ID:       "9897f4c66b5e",
							Created:  time.Now().Unix(),
							RepoTags: []string{"<none>:<none>"},
							Labels:   map[string]string{},
						},
					}, nil).
					Times(1)
			},
			wantErr: nopErr,
			expected: []docker.Image{
				{
					ID:     "a76d6a1f0270",
					Tags:   []string{"nginx:latest"},
					Labels: map[string]string{},
				},
				{
					ID:     "9897f4c66b5e",
					Tags:   []string{"<none>:<none>"},
					Labels: map[string]string{},
				},
			},
		},
		{
			name: "filter images by label",
			args: args{
				ctx: context.Background(),
				options: docker.ExpiredImageListOptions{
					LifetimeThresholdInDays: 50,
					IgnoreLabels:            []string{"com.github.lucasmendesl.beerus.testLabel"},
				},
			},
			mockSetup: func() {
				dockerClient.
					EXPECT().
					ImageList(
						gomock.Any(),
						gomock.Any(),
					).
					Return([]image.Summary{
						{
							ID:       "104340b97284",
							Created:  time.Now().Add(-time.Hour * 24 * 55).Unix(),
							RepoTags: []string{"beerus:latest"},
							Labels:   map[string]string{"com.github.lucasmendesl.beerus.service": "true"},
						},
						{
							ID:       "492084a114c1",
							Created:  time.Now().Add(-time.Hour * 24 * 70).Unix(),
							RepoTags: []string{"nginx:latest"},
							Labels:   map[string]string{"com.github.lucasmendesl.beerus.testLabel": "true"},
						},
						{
							ID:       "b320553669f9",
							Created:  time.Now().Add(-time.Hour * 24 * 55).Unix(),
							RepoTags: []string{"php:latest"},
							Labels:   map[string]string{},
						},
					}, nil).
					Times(1)
			},
			wantErr: nopErr,
			expected: []docker.Image{
				{
					ID:     "b320553669f9",
					Tags:   []string{"php:latest"},
					Labels: map[string]string{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			d := docker.New(dockerClient, logger)

			got, err := d.ListExpiredImages(tt.args.ctx, tt.args.options)
			if tt.wantErr(t, err) {
				return
			}

			require.Equal(t, tt.expected, got)
		})
	}
}
