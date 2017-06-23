package app

import (
	"github.com/docker/docker/client"
	"killl/lib/httpcli"
	"strings"
	"time"
)

type DockerCli struct {
	client client.APIClient
}

func NewDockerCli(config *Config) (DockerCli, error) {
	if strings.HasPrefix(config.dockerHost, "unix:///") {
		dockerCli, err := client.NewClient(config.dockerHost, config.dockerAPIVersion, nil, nil)
		if err != nil {
			return DockerCli{}, err
		}
		return DockerCli{client: dockerCli}, nil
	}
	httpCli := httpcli.NewHttpClient(
		httpcli.Timeout(time.Duration(config.requestTimeout)),
		httpcli.KeepAlive(time.Duration(config.keepAlive)),
	)
	dockerCli, err := client.NewClient(config.dockerHost,
		config.dockerAPIVersion,
		httpCli,
		nil)
	if err != nil {
		return DockerCli{}, err
	}
	return DockerCli{client: dockerCli}, nil
}

func (cli *DockerCli) Client() client.APIClient {
	return cli.client
}
