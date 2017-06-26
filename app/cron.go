package app

import (
	"context"
	"github.com/docker/docker/api/types"
	"killl/lib/log"
	"killl/lib/ring"
	"strings"
	"time"
)

func CronCheckContainers(ctx context.Context, interval time.Duration, cli DockerCli) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			log.Debug("cron check containers")
			if err := checkContainers(ctx, cli); err != nil {
				log.Errorf("check container error: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func checkContainers(ctx context.Context, cli DockerCli) error {
	var (
		name      string
		id        string
		clearTemp = make(map[string]struct{})
	)
	containers, err := cli.Client().ContainerList(ctx, types.ContainerListOptions{Quiet: false, All: true, Size: false})
	if err != nil {
		return err
	}
	for _, container := range containers {
		if len(container.Names) < 1 {
			log.Errorf("invalid container name: %s", container.Names)
			continue
		}
		name = container.Names[0]
		if strings.HasPrefix(name, "/") {
			name = name[1:]
		}
		id = container.ID
		if !validContainer(name, id) {
			if err := HandleContainer(ctx, cli, &ring.Payload{Name: name, ID: id}); err != nil {
				log.Errorf("handle container error: %v", err)
				continue
			}
		}
		clearTemp[id] = struct{}{}
	}
	// clear cache
	ClearCache(clearTemp)
	return nil
}

func validContainer(name string, id string) bool {
	if ExistInWhiteList(name) || ExistInContainerCache(id) || HasSpecificPrefix(name) {
		return true
	}
	return false
}

func CronImageClean(ctx context.Context, interval, before time.Duration, cli DockerCli) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			log.Debug("cron image clean")
			if err := cleanImage(ctx, before, cli); err != nil {
				log.Errorf("clean image error: %v", err)
			}
		}
	}
}

func cleanImage(ctx context.Context, before time.Duration, cli DockerCli) error {
	images, err := cli.Client().ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return err
	}
	for _, image := range images {
		tm := time.Unix(image.Created, 0)
		if tm.Add(before).Before(time.Now()) {
			log.Infof("remove image: id(%s)", image.ID)
			if _, err := cli.Client().ImageRemove(ctx, image.ID, types.ImageRemoveOptions{}); err != nil {
				log.Errorf("remove image error: %v", err)
			}
		}
	}
	return nil
}
