package docker

import (
	"context"
	"github.com/docker/docker/api/types/filters"
	"io"
	"time"

	"github.com/docker/docker/api/types"
)

// RemoveContainerByID 根据容器ID删除容器
func RemoveContainerByID(containerID string) error { return c.RemoveContainerByID(containerID) }

// RemoveContainerByID 根据容器ID删除容器
func (c *ClientEx) RemoveContainerByID(containerID string) error {
	return c.ContainerRemove(c.ctx,
		containerID,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   true,
			Force:         true,
		})
}

// ListContainer 列出容器
//     all：是否列出所有容器
func ListContainer(all ...bool) ([]types.Container, error) { return c.ListContainer(all...) }

// 清除docker镜像
func ClearDockerContainer() error {
	return c.ClearDockerContainer()
}

// 查询容器的详情
func InspectByContainerName(containerName string) (types.ContainerJSON, error) {
	return c.InspectByContainerName(containerName)
}

// 登录
func Login(userName string, password string) error {
	return c.login(userName, password)
}

// 通过时间查询container log
func QueryContainerLogsByDate(containerId string, startDay string, endDay string) (io.ReadCloser, error) {
	return c.QueryContainerLogsByDate(containerId, startDay, endDay)
}

// ListContainer 列出容器
//     all：是否列出所有容器
func (c *ClientEx) ListContainer(all ...bool) ([]types.Container, error) {
	listAll := append(all, false)[0]
	return c.ContainerList(c.ctx, types.ContainerListOptions{All: listAll})
}

// StartContainer 启动容器
func StartContainer(containerID string) error { return c.StartContainer(containerID) }

// StartContainer 启动容器
func (c *ClientEx) StartContainer(containerID string) error {
	return c.ContainerStart(c.ctx, containerID, types.ContainerStartOptions{})
}

// StopContainer 停止容器
func StopContainer(containerID string) error { return c.StopContainer(containerID) }

// StopContainer 停止容器
func (c *ClientEx) StopContainer(containerID string) error {
	timeout := 10 * time.Second
	return c.ContainerStop(c.ctx, containerID, &timeout)
}

// RestartContainer 重启容器
func RestartContainer(containerID string) error { return c.RestartContainer(containerID) }

// RestartContainer 重启容器
func (c *ClientEx) RestartContainer(containerID string) error {
	timeout := 10 * time.Second
	return c.ContainerRestart(c.ctx, containerID, &timeout)
}

func (c *ClientEx) ClearDockerContainer() error {
	ctx := context.Background()
	_, e := c.ContainersPrune(ctx, filters.Args{})
	//s:=pruneReport.ContainersDeleted
	//klog.Info(s)
	return e
}

func (c *ClientEx) InspectByContainerName(containerName string) (types.ContainerJSON, error) {
	return c.ContainerInspect(context.Background(), containerName)
}

func (c *ClientEx) login(userName string, password string) error {
	config := types.AuthConfig{
		Username: userName,
		Password: password,
	}
	_, err := c.RegistryLogin(context.Background(), config)
	return err
}

func (c *ClientEx) QueryContainerLogsByDate(containerId string, startDay string, endDay string) (io.ReadCloser, error) {
	return c.ContainerLogs(context.Background(), containerId, types.ContainerLogsOptions{
		ShowStdout: true,
		Since:      startDay,
		Until:      endDay,
	})
}
