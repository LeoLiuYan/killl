package app

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"killl/lib/log"
	ring2 "killl/lib/ring"
)

const (
	SIGNAL = "KILL"
)

func RunEvents(ctx context.Context, dockercli DockerCli, opts types.EventsOptions) (<-chan events.Message, <-chan error) {
	return dockercli.Client().Events(ctx, opts)
}

func HandleEvent(event events.Message, ring *ring2.Ring) error {
	if event.Action == "create" &&
		!ExistInWhiteList(event.Actor.Attributes["name"]) &&
		!HasSpecificPrefix(event.Actor.Attributes["name"]) {
		payload, err := ring.Set()
		if err != nil {
			return err
		}
		payload.ID = event.ID
		payload.Name = event.Actor.Attributes["name"]
		ring.SetA()
	}
	return nil
}

func HandleContainer(ctx context.Context, dockerCli DockerCli, container *ring2.Payload) error {
	containerJson, err := dockerCli.Client().ContainerInspect(ctx, container.ID)
	if err != nil {
		return err
	}
	for k, v := range containerJson.Config.Labels {
		if ExistLabelsList(k, v) {
			// cache
			log.Infof("cache container: id(%s) name(%s)", container.ID, container.Name)
			ContainerCache.put(container.ID)
			return nil
		}
	}
	log.Infof("kill container: id(%s) name(%s)", container.ID, container.Name)
	if err := dockerCli.Client().ContainerKill(ctx, container.ID, SIGNAL); err != nil {
		log.Errorf("kill container: id(%s) error", container.ID)
	}
	return dockerCli.Client().ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{false, false, true})
}
