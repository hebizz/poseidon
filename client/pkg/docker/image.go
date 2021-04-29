package docker

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	log "k8s.io/klog"
)

// BuildImage build image
func BuildImage(buildContext io.Reader, imageName string) error {
	return c.BuildImage(buildContext, imageName)
}

// BuildImage build image
func (c *ClientEx) BuildImage(buildContext io.Reader, imageName string) error {
	buildOptions := types.ImageBuildOptions{
		Tags: []string{imageName},
	}

	output, err := c.ImageBuild(c.ctx, buildContext, buildOptions)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(body), "error") {
		return fmt.Errorf("build image to docker error")
	}

	return nil
}

// ListImage list docker images
func ListImage() ([]types.ImageSummary, error) { return c.ListImage() }

// ListImage list docker images
func (c *ClientEx) ListImage() ([]types.ImageSummary, error) {
	return c.ImageList(c.ctx, types.ImageListOptions{})
}

// RemoveAllImage 删除所有docker image
func RemoveAllImage() error { return c.RemoveAllImage() }

// RemoveAllImage 删除所有docker image
func (c *ClientEx) RemoveAllImage() error {
	images, err := c.ImageList(c.ctx, types.ImageListOptions{})
	if err != nil {
		log.Error("list images error: %v", err)
		return err
	}

	for _, image := range images {
		_, err := c.ImageRemove(c.ctx, image.ID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
		if err != nil {
			log.Error(err)
		}
		log.Info("has delete image : " + image.ID)
	}

	return nil
}

// LoadImage load image
func LoadImage(input io.Reader) error { return c.LoadImage(input) }

// LoadImage load image
func (c *ClientEx) LoadImage(input io.Reader) error {
	output, err := c.ImageLoad(c.ctx, input, true)
	log.Infof("load image output: %+v", output)
	return err
}
