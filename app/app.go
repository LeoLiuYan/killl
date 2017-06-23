package app

import (
	"context"
	"github.com/docker/docker/api/types"
	"io"
	"killl/lib/log"
	ring2 "killl/lib/ring"
	"strings"
	"time"
)

const (
	CONSUMER_INTERVAL = 1000 * time.Millisecond
)

var (
	Whitelist = make(map[string]struct{})
	// todo: support label Selectors as Kubernetes
	LabelsList     = make(map[string]string)
	ContainerCache = make(map[string]struct{})
	mesosPrefix string
)

func Run(ctx context.Context, config *Config) {
	log.Debugf("config: %v", config)
	mesosPrefix = config.mesosPrefix
	dockerCli, err := NewDockerCli(config)
	if err != nil {
		log.Fatalf("Run() error: %v", err)
	}
	// connect to daemon
	_, err = dockerCli.Client().ServerVersion(ctx)
	if err != nil {
		log.Fatalf("connect to daemon error: %v", err)
	}
	initWhiteList(config)
	initLabelsList(config)
	// ring buffer
	ring := ring2.NewRing(config.ringBuffer)
	go consumer(ctx, dockerCli, ring)
	eventOpts := types.EventsOptions{}
	events, errs := RunEvents(ctx, dockerCli, eventOpts)
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errs:
			if err != io.EOF {
				log.Errorf("receive events error: %v", err)
			}
		case event := <-events:
			if err := HandleEvent(event, ring); err != nil {
				log.Errorf("handler event error: %v", err)
			}
		}
	}
}

func consumer(ctx context.Context, dockerCli DockerCli, ring *ring2.Ring) {
	ticker := time.NewTicker(CONSUMER_INTERVAL)
loop:
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c, err := ring.Get()
			if err != nil {
				if err != ring2.ErrRingEmpty {
					log.Errorf("consumer event error: %v", err)
				}
				continue loop
			}
			if err := HandleContainer(ctx, dockerCli, c); err != nil {
				log.Errorf("handle container(%s) error: %v", c, err)
			}
			ring.GetA()
		}
	}
}

func initLabelsList(config *Config) {
	var key, value string
	for _, label := range config.labels {
		ll := strings.Split(label, "=")
		if len(ll) < 2 {
			log.Errorf("label invalid format: %v", label)
			continue
		}
		key, value = ll[0], ll[1]
		LabelsList[key] = value
	}
}

func ExistLabelsList(key, value string) bool {
	if key == "" {
		return false
	}
	if v, ok := LabelsList[key]; !ok || v != value {
		return false
	}
	return true
}

func initWhiteList(config *Config) {
	for _, name := range config.whitelist {
		Whitelist[name] = struct{}{}
	}
}

func ExistInWhiteList(name string) bool {
	if name == "" {
		return false
	}
	if _, ok := Whitelist[name]; !ok {
		return false
	}
	return true
}

func ExistInContainerCache(containerID string) bool {
	if _, ok := ContainerCache[containerID]; !ok {
		return false
	}
	return true
}

func IsMesosPrefix(name string) bool {
	return strings.HasPrefix(name, mesosPrefix)
}
