package container

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"
)

type MockContainerUpdate func(container.InspectResponse, image.InspectResponse)

func MockContainer(updates ...MockContainerUpdate) *Container {
	containerInfo := container.InspectResponse{
		ContainerJSONBase: &container.ContainerJSONBase{
			ID:         "container_id",
			Image:      "image",
			Name:       "test-containrrr",
			HostConfig: &container.HostConfig{},
		},
		Config: &container.Config{
			Labels: map[string]string{},
		},
	}
	image := image.InspectResponse{
		ID: "image_id",
	}

	for _, update := range updates {
		update(containerInfo, image)
	}
	return NewContainer(&containerInfo, &image)
}

func WithPortBindings(portBindingSources ...string) MockContainerUpdate {
	return func(cnt container.InspectResponse, img image.InspectResponse) {
		portBindings := nat.PortMap{}
		for _, pbs := range portBindingSources {
			portBindings[nat.Port(pbs)] = []nat.PortBinding{}
		}
		cnt.HostConfig.PortBindings = portBindings
	}
}

func WithImageName(name string) MockContainerUpdate {
	return func(cnt container.InspectResponse, img image.InspectResponse) {
		cnt.Config.Image = name
		img.RepoTags = append(img.RepoTags, name)
	}
}

func WithLinks(links []string) MockContainerUpdate {
	return func(cnt container.InspectResponse, img image.InspectResponse) {
		cnt.HostConfig.Links = links
	}
}

func WithLabels(labels map[string]string) MockContainerUpdate {
	return func(cnt container.InspectResponse, img image.InspectResponse) {
		cnt.Config.Labels = labels
	}
}

func WithContainerState(state types.ContainerState) MockContainerUpdate {
	return func(cnt container.InspectResponse, img image.InspectResponse) {
		cnt.State = &container.State{
			Status: state.Status,
		}
	}
}

func WithHealthcheck(healthConfig container.HealthConfig) MockContainerUpdate {
	return func(cnt container.InspectResponse, img image.InspectResponse) {
		cnt.Config.Healthcheck = &healthConfig
	}
}

func WithImageHealthcheck(healthConfig container.HealthConfig) MockContainerUpdate {
	return func(cnt container.InspectResponse, img image.InspectResponse) {
		img.Config.Healthcheck = &healthConfig
	}
}
