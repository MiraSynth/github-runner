package container

import (
	"context"
	"io"
	"os"
	"runtime"

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

type Container interface {
	Create(context.Context, *Options) (string, error)
	Start(context.Context, string) error
}

type implementation struct {
	client *client.Client
}

func New() (Container, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	impl := &implementation{
		client: c,
	}

	runtime.SetFinalizer(impl, func(impl implementation) {
		if impl.client == nil {
			return
		}

		defer impl.client.Close()
	})

	return impl, nil
}

func (c *implementation) Create(ctx context.Context, options *Options) (string, error) {
	reader, err := c.client.ImagePull(ctx, options.ImageName, image.PullOptions{})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return "", err
	}

	containerConfig := &container.Config{
		Image:      options.ImageName,
		Cmd:        options.Command,
		Entrypoint: options.Entrypoint,
		Env:        options.Environment,
	}
	createResponse, err := c.client.ContainerCreate(ctx, containerConfig, nil, nil, nil, options.Name)
	if err != nil {
		return "", err
	}

	return createResponse.ID, nil
}

func (c *implementation) Start(ctx context.Context, containerId string) error {
	containerStartOptions := container.StartOptions{}
	err := c.client.ContainerStart(ctx, containerId, containerStartOptions)
	if err != nil {
		return err
	}

	containerLogsOptions := container.LogsOptions{ShowStdout: true, Follow: true}
	out, err := c.client.ContainerLogs(ctx, containerId, containerLogsOptions)
	if err != nil {
		return err
	}

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		return err
	}

	return nil
}
