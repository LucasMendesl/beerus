package cleaner_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/lucasmendesl/beerus/cleaner"
	"github.com/lucasmendesl/beerus/config"
	"github.com/lucasmendesl/beerus/docker"
	mock "github.com/lucasmendesl/beerus/docker/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type wantErr func(t *testing.T, err error) bool

func TestCleaner_Run(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)

		config = &config.Beerus{
			ConcurrencyLevel:        1,
			ExpirePollCheckInterval: 1,
			Images: config.Image{
				LifetimeThreshold: 1,
			},
			Containers: config.Container{
				MaxAlwaysRestartPolicyCount: 1,
			},
		}

		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(dockerAPI *mock.MockBeerusContainerAPI)
		wantErr   wantErr
	}{
		{
			name: "error listing containers",
			args: args{
				ctx: context.Background(),
			},
			setupMock: func(dockerAPI *mock.MockBeerusContainerAPI) {
				dockerAPI.
					EXPECT().
					ListContainers(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return(nil, errors.New("error listing containers")).
					AnyTimes()
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "error listing containers")
				return true
			},
		},
		{
			name: "error removing containers",
			args: args{
				ctx: context.Background(),
			},
			setupMock: func(dockerAPI *mock.MockBeerusContainerAPI) {
				dockerAPI.
					EXPECT().
					ListContainers(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]docker.Container{
						{
							ID:           "cadc6990a82e",
							Image:        "nginx:1.27.3-alpine",
							ImageID:      "sha256:b4ef436c698b07",
							Labels:       map[string]string{},
							CreatedAt:    time.Time{},
							Status:       "stopped",
							RestartCount: 0,
							RestartPolicy: container.RestartPolicy{
								Name: "no",
							},
						},
					},
						nil).
					AnyTimes()

				dockerAPI.
					EXPECT().
					RemoveContainer(
						gomock.Any(),
						gomock.Any(),
					).
					Return(errors.New("error removing containers")).
					AnyTimes()
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "error removing container with id cadc6990a82e: error removing containers")
				return true
			},
		},
		{
			name: "error listing images",
			args: args{
				ctx: context.Background(),
			},
			setupMock: func(dockerAPI *mock.MockBeerusContainerAPI) {
				dockerAPI.
					EXPECT().
					ListContainers(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]docker.Container{
						{
							ID:           "cadc6990a82e",
							Image:        "nginx:1.27.3-alpine",
							ImageID:      "sha256:b4ef436c698b07",
							Labels:       map[string]string{},
							CreatedAt:    time.Time{},
							Status:       "stopped",
							RestartCount: 0,
							RestartPolicy: container.RestartPolicy{
								Name: "no",
							},
						},
					},
						nil).
					AnyTimes()

				dockerAPI.
					EXPECT().
					RemoveContainer(
						gomock.Any(),
						gomock.Any(),
					).
					AnyTimes()

				dockerAPI.
					EXPECT().
					ListExpiredImages(
						gomock.Any(),
						gomock.Any(),
					).
					Return(nil, errors.New("error listing images")).
					AnyTimes()
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "error getting expired images: error listing images")
				return true
			},
		},
		{
			name: "error removing images",
			args: args{
				ctx: context.Background(),
			},
			setupMock: func(dockerAPI *mock.MockBeerusContainerAPI) {
				dockerAPI.
					EXPECT().
					ListContainers(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]docker.Container{
						{
							ID:           "cadc6990a82e",
							Image:        "nginx:1.27.3-alpine",
							ImageID:      "sha256:b4ef436c698b07",
							Labels:       map[string]string{},
							CreatedAt:    time.Time{},
							Status:       "stopped",
							RestartCount: 0,
							RestartPolicy: container.RestartPolicy{
								Name: "no",
							},
						},
					},
						nil).
					AnyTimes()

				dockerAPI.
					EXPECT().
					RemoveContainer(
						gomock.Any(),
						gomock.Any(),
					).
					AnyTimes()

				dockerAPI.
					EXPECT().
					ListExpiredImages(
						gomock.Any(),
						gomock.Any(),
					).
					Return([]docker.Image{
						{
							ID:     "b0757c55a1fd",
							Labels: map[string]string{},
							Tags:   []string{"docker:stable"},
						},
					}, nil).
					AnyTimes()

				dockerAPI.
					EXPECT().
					RemoveImage(
						gomock.Any(),
						gomock.Any(),
					).
					Return(errors.New("error removing images")).
					AnyTimes()
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "error removing image with id b0757c55a1fd: error removing images")
				return true
			},
		},
		{
			name: "clean resources",
			args: args{
				ctx: context.Background(),
			},
			setupMock: func(dockerAPI *mock.MockBeerusContainerAPI) {
				dockerAPI.
					EXPECT().
					ListContainers(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]docker.Container{
						{
							ID:           "cadc6990a82e",
							Image:        "nginx:1.27.3-alpine",
							ImageID:      "sha256:b4ef436c698b07",
							Labels:       map[string]string{},
							CreatedAt:    time.Time{},
							Status:       "stopped",
							RestartCount: 0,
							RestartPolicy: container.RestartPolicy{
								Name: "no",
							},
						},
					},
						nil).
					AnyTimes()

				dockerAPI.
					EXPECT().
					RemoveContainer(
						gomock.Any(),
						gomock.Any(),
					).
					AnyTimes()

				dockerAPI.
					EXPECT().
					ListExpiredImages(
						gomock.Any(),
						gomock.Any(),
					).
					Return([]docker.Image{
						{
							ID:     "b0757c55a1fd",
							Labels: map[string]string{},
							Tags:   []string{"docker:stable"},
						},
					}, nil).
					AnyTimes()

				dockerAPI.
					EXPECT().
					RemoveImage(
						gomock.Any(),
						gomock.Any(),
					).
					Return(errors.New("error removing images")).
					AnyTimes()
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "error removing image with id b0757c55a1fd: error removing images")
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dockerAPI := mock.NewMockBeerusContainerAPI(ctrl)

			tt.setupMock(dockerAPI)
			cleaner := cleaner.New(dockerAPI, config, logger)

			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				err := cleaner.Run(ctx)

				if tt.wantErr(t, err) {
					return
				}
				cancel()
			}()

			wg.Wait()
		})
	}
}
