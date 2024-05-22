package container

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Options struct {
	Name        string   `json:"string"`
	ImageName   string   `json:"imageName"`
	Command     []string `json:"command"`
	Entrypoint  []string `json:"entrypoint"`
	Environment []string `json:"environment"`
}

func Start(options *Options) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, options.ImageName, image.PullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	containerConfig := &container.Config{
		Image:      options.ImageName,
		Cmd:        options.Command,
		Entrypoint: options.Entrypoint,
		Env:        options.Environment,
	}
	createResponse, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, options.Name)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, createResponse.ID, container.StartOptions{}); err != nil {
		return err
	}

	out, err := cli.ContainerLogs(ctx, createResponse.ID, container.LogsOptions{ShowStdout: true, Follow: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}
