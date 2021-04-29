package docker

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
)

// ClientEx docker client 的扩展，方便调用
type ClientEx struct {
	*client.Client

	ctx context.Context
}

var c *ClientEx

func init() {
	c, _ = NewClientExWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// NewClientEx initializes a new API client for the given host and API version
func NewClientEx(host string, version string, httpClient *http.Client, httpHeaders map[string]string) (*ClientEx, error) {
	dockerClient, err := client.NewClient(host, version, httpClient, httpHeaders)
	if err != nil {
		return nil, err
	}

	return &ClientEx{
		Client: dockerClient,
		ctx:    context.Background(),
	}, nil
}

// NewEnvCLientEx initializes a new API client based on environment variables
func NewEnvCLientEx() (*ClientEx, error) {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &ClientEx{
		Client: dockerClient,
		ctx:    context.Background(),
	}, nil
}

// NewClientExWithOpts initializes a new API client with default values
func NewClientExWithOpts(opt ...client.Opt) (*ClientEx, error) {
	dockerClient, err := client.NewClientWithOpts(opt...)
	if err != nil {
		return nil, err
	}

	return &ClientEx{
		Client: dockerClient,
		ctx:    context.Background(),
	}, nil
}
