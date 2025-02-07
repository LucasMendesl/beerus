package docker

import (
	"context"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

const (
	eventAction = "event"
	filterType  = "type"
)

// FromEvents returns a channel of EventResult objects, which contain either an
// events.Message struct describing a Docker event, or an error if there is an
// issue fetching the event stream.
//
// The function takes a context.Context and a variadic list of
// events.Action values, which are used to filter the types of events that are
// returned.
//
// The function returns a channel of EventResult objects, which is closed when
// the context is canceled or when an error occurs. If an error occurs, the
// channel will contain a single EventResult object with a non-nil Err field.
// If the context is canceled, the channel will contain a single EventResult
// object with an Err field that is equal to the context's error.
func (d *dockerClient) FromEvents(ctx context.Context, actions ...events.Action) <-chan EventResult {
	filterOpts := filters.NewArgs(
		filters.Arg(filterType, string(events.ContainerEventType)),
		filters.Arg(filterType, string(events.ImageEventType)),
	)
	for _, action := range actions {
		filterOpts.Add(eventAction, string(action))
	}

	options := events.ListOptions{
		Filters: filterOpts,
	}

	eventCh := make(chan EventResult, 1)
	go func() {
		msgCh, errCh := d.cli.Events(ctx, options)
		defer close(eventCh)

		for {
			select {
			case <-ctx.Done():
				eventCh <- EventResult{Err: ctx.Err()}
				return
			case msg := <-msgCh:
				eventCh <- EventResult{Message: msg}
			case err := <-errCh:
				eventCh <- EventResult{Err: err}
				return
			}
		}
	}()

	return eventCh
}
