package docker_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/lucasmendesl/beerus/docker"
	mock "github.com/lucasmendesl/beerus/docker/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDockerClient_ListContainers(t *testing.T) {
	var (
		ctrl         = gomock.NewController(t)
		dockerClient = mock.NewMockClient(ctrl)
		logger       = slog.New(slog.NewJSONHandler(io.Discard, nil))

		createdAt              = time.Date(2025, time.January, 4, 12, 0, 0, 0, time.Local)
		containerListError     = errors.New("container list error")
	)
	type args struct {
		ctx         context.Context
		concurrency uint8
		options     []docker.ListContainersOptions
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		expected  []docker.Container
		wantErr   wantErr
	}{
		{
			name: "list container error",
			args: args{
				ctx:         context.Background(),
				concurrency: 1,
				options:     []docker.ListContainersOptions{},
			},
			mockSetup: func() {
				dockerClient.
					EXPECT().
					ContainerList(
						gomock.Any(),
						gomock.Any(),
					).
					Return(nil, containerListError).
					Times(1)
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "fetching containers error: container list error")
				return true
			},
		},
		{
			name: "list container empty",
			args: args{
				ctx:         context.Background(),
				concurrency: 1,
				options:     []docker.ListContainersOptions{},
			},
			mockSetup: func() {
				dockerClient.
					EXPECT().
					ContainerList(
						gomock.Any(),
						gomock.Any(),
					).
					Return([]types.Container{}, nil).
					Times(1)
			},
			expected: []docker.Container{},
			wantErr:  nopErr,
		},
		{
			name: "filter containers by labels",
			args: args{
				ctx:         context.Background(),
				concurrency: 1,
				options: []docker.ListContainersOptions{
					docker.WithContainerLabel("com.github.lucasmendesl.beerus.testLabel"),
				},
			},
			mockSetup: func() {
				dockerClient.
					EXPECT().
					ContainerList(
						gomock.Any(),
						gomock.Any(),
					).
					Return([]types.Container{
						{
							ID:      "d0fcf186fa",
							Image:   "beerus:latest",
							Created: createdAt.Unix(),
							Labels:  map[string]string{"com.github.lucasmendesl.beerus.service": "true"},
							Status:  "stopped",
							ImageID: "sha256:d55c68fb34057c75d9f0",
						},
						{
							ID:      "b0757c55a1fd",
							Image:   "busybox:latest",
							Created: createdAt.Unix(),
							Labels:  map[string]string{},
							Status:  "stopped",
							ImageID: "sha256:b5ad7243b38d33a8db255",
						},
						{
							ID:      "b4ef436c698",
							Image:   "nginx:latest",
							Created: createdAt.Unix(),
							Labels:  map[string]string{"com.github.lucasmendesl.beerus.testLabel": "true"},
							Status:  "stopped",
							ImageID: "sha256:d03820ba684b9a7ce9b13",
						},
					},
						nil).
					Times(1)
			},
			expected: []docker.Container{
				{
					ID:           "b0757c55a1fd",
					Image:        "busybox:latest",
					ImageID:      "sha256:b5ad7243b38d33a8db255",
					Labels:       map[string]string{},
					Status:       "stopped",
					CreatedAt:    createdAt,
				},
			},
			wantErr: nopErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			d := docker.New(dockerClient, logger)

			got, err := d.ListContainers(tt.args.ctx, tt.args.concurrency, tt.args.options...)
			if tt.wantErr(t, err) {
				return
			}

			require.Equal(t, tt.expected, got)
		})
	}
}
