package cleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/lucasmendesl/beerus/docker"
)

// watch sets up event listeners for specific Docker events and logs them. It
// takes a context.Context and a channel of error objects as parameters. The
// function listens for specific Docker events, such as image untagging and
// container exit events, and logs these events. It also runs a periodic task to
// identify and log removable images based on certain criteria. The function
// returns an error if any occurs during the cleanup process.
func (c *cleaner) watch(ctx context.Context, errCh chan<- error) {
	// run the image checker periodically, following the configuration
	go c.pollImageChecker(ctx, errCh)

	c.log.Info("Starting watching docker events...", "context", "Event")
	// listen for specific Docker events
	// container exit events
	// image untagging events
	for result := range c.d.FromEvents(ctx,
		events.ActionDie,
		events.ActionUnTag,
	) {
		if result.Err != nil {
			c.log.Error("error receiving event", "error", result.Err, "context", "Event")
			errCh <- result.Err
			return
		}

		c.log.Debug("event received", "action", result.Message.Action, "id", result.Message.ID, "context", "Event")
		c.handleWatcherEvent(ctx, result.Message)
	}
}

// pollImageChecker is a goroutine that periodically checks for removable
// images and removes them. It takes a context.Context, a cleaner object, and a
// channel of error objects as parameters. The function runs in an infinite
// loop, checking for removable images once a minute. If an error occurs during
// the cleanup process, the function sends the error on the error channel and
// returns.
func (c *cleaner) pollImageChecker(ctx context.Context, errCh chan<- error) {
	c.log.Info("Starting periodic image checker, checking for removable images every", "interval in hours", c.config.ExpirePollCheckInterval, "context", "Image Poller")

	ticker := time.NewTicker(time.Hour * time.Duration(c.config.ExpirePollCheckInterval))
	defer ticker.Stop()

	for range ticker.C {
		c.log.Debug("Checking for removable images", "context", "Image Poller")
		removableImgs, err := c.listAllowedImagesToRemove(ctx)

		if err != nil {
			errCh <- fmt.Errorf("list image poller error: %w", err)
			return
		}

		c.log.Debug("Found removable images", "count", len(removableImgs), "context", "Image Poller")
		if err := c.removeImages(ctx, removableImgs...); err != nil {
			errCh <- fmt.Errorf("remove image poller error: %w", err)
			return
		}
	}
}

// handleWatcherEvent takes a context.Context and a Docker events.Message object as parameters.
// The function inspects the message's Action field to determine how to handle the event.
// If the action is "untag", the function removes the image if it is not used by any containers.
// If the action is "die", the function inspects the container that exited and removes it if it does
// not have a restart policy.
func (c *cleaner) handleWatcherEvent(ctx context.Context, message events.Message) {
	switch message.Action {
	case events.ActionUnTag:
		// if an image is untagged, remove it if it is not used by any
		// containers
		c.log.Debug("untag event received, removing image", "id", message.ID, "context", "Event")
		if err := c.removeImages(ctx, docker.Image{ID: message.ID}); err != nil {
			c.log.Error("error on removing image", "context", "Event", "err", err)
		}
	case events.ActionDie:
		// if a container exits, remove it if it does not have a restart
		// policy
		c.log.Debug("die event received, inspecting container", "id", message.ID, "context", "Event")
		containerDetails, err := c.d.Inspect(ctx, message.ID)
		if err != nil {
			c.log.Error("error inspecting container", "error", err, "context", "Event")
			break
		}

		c.log.Debug("container inspected", "id", message.ID, "restart-policy", containerDetails.HostConfig.RestartPolicy.Name, "context", "Event")

		container := docker.Container{
			RestartCount:  containerDetails.RestartCount,
			RestartPolicy: containerDetails.HostConfig.RestartPolicy,
		}
		if !docker.CanRemoveContainer(container, c.config.Containers.MaxAlwaysRestartPolicyCount) {
			c.log.Debug("unavailable container to remove", "id", message.ID, "restart-policy", containerDetails.HostConfig.RestartPolicy.Name, "context", "Event")
			break
		}

		c.log.Debug("container is removable, removing it", "id", message.ID, "context", "Event")
		if err := c.removeContainers(ctx, message.Actor.ID); err != nil {
			c.log.Error("removing container", "context", "Event", "err", err)
		}
	}
}
