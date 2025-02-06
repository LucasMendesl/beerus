package docker_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"

	"github.com/docker/docker/api/types/events"
	"github.com/lucasmendesl/beerus/docker"
	mock "github.com/lucasmendesl/beerus/docker/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDockerClient_FromEvents(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)

		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	)
	type args struct {
		ctx    context.Context
		action events.Action
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(mockCli *mock.MockClient)
		expected  events.Message
		wantErr   wantErr
	}{
		{
			name: "handle successful event",
			args: args{
				ctx:    context.Background(),
				action: events.ActionDie,
			},
			setupMock: func(dockerApi *mock.MockClient) {
				eventCh := make(chan events.Message, 1)
				errCh := make(chan error, 1)
				dockerApi.
					EXPECT().
					Events(
						gomock.Any(),
						gomock.Any(),
					).Return(eventCh, errCh).
					Times(1)
				eventCh <- events.Message{Type: "container", Action: events.ActionDie}
			},
			expected: events.Message{Type: "container", Action: events.ActionDie},
			wantErr:  nopErr,
		},
		{
			name: "error fetching events",
			setupMock: func(dockerApi *mock.MockClient) {
				eventCh := make(chan events.Message, 1)
				errCh := make(chan error, 1)
				dockerApi.
					EXPECT().
					Events(
						gomock.Any(),
						gomock.Any(),
					).
					Return(eventCh, errCh).
					Times(1)
				errCh <- errors.New("error fetching events")
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "error fetching events")
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCli := mock.NewMockClient(ctrl)
			dockerClient := docker.New(mockCli, logger)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			tt.setupMock(mockCli)
			resultCh := dockerClient.FromEvents(ctx, events.ActionCreate)

			go func() {
				defer wg.Done()
				got := <-resultCh

				if tt.wantErr(t, got.Err) {
					return
				}

				require.Equal(t, tt.expected, got.Message)
				cancel()
			}()

			wg.Wait()
		})
	}
}
