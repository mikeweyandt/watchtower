package container

import (
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"
	dockerspec "github.com/moby/docker-image-spec/specs-go/v1"
)

type MockContainerUpdate func(dockerContainer.InspectResponse, image.InspectResponse)

func MockContainer(updates ...MockContainerUpdate) *Container {
	containerInfo := dockerContainer.InspectResponse{
		ContainerJSONBase: &dockerContainer.ContainerJSONBase{
			ID:         "container_id",
			Image:      "image",
			Name:       "test-containrrr",
			HostConfig: &dockerContainer.HostConfig{},
		},
		Config: &dockerContainer.Config{
			Labels: map[string]string{},
		},
	}
	image := image.InspectResponse{
		ID:     "image_id",
		Config: &dockerspec.DockerOCIImageConfig{},
	}

	for _, update := range updates {
		update(containerInfo, image)
	}
	return NewContainer(&containerInfo, &image)
}

func WithPortBindings(portBindingSources ...string) MockContainerUpdate {
	return func(cnt dockerContainer.InspectResponse, img image.InspectResponse) {
		portBindings := nat.PortMap{}
		for _, pbs := range portBindingSources {
			portBindings[nat.Port(pbs)] = []nat.PortBinding{}
		}
		cnt.HostConfig.PortBindings = portBindings
	}
}

func WithImageName(name string) MockContainerUpdate {
	return func(cnt dockerContainer.InspectResponse, img image.InspectResponse) {
		cnt.Config.Image = name
		img.RepoTags = append(img.RepoTags, name)
	}
}

func WithLinks(links []string) MockContainerUpdate {
	return func(cnt dockerContainer.InspectResponse, img image.InspectResponse) {
		cnt.HostConfig.Links = links
	}
}

func WithLabels(labels map[string]string) MockContainerUpdate {
	return func(cnt dockerContainer.InspectResponse, img image.InspectResponse) {
		cnt.Config.Labels = labels
	}
}

func WithContainerState(state dockerContainer.State) MockContainerUpdate {
	return func(cnt dockerContainer.InspectResponse, img image.InspectResponse) {
		cnt.State = &dockerContainer.State{
			Status: state.Status,
		}
	}
}

func WithHealthcheck(healthConfig dockerContainer.HealthConfig) MockContainerUpdate {
	return func(cnt dockerContainer.InspectResponse, img image.InspectResponse) {
		cnt.Config.Healthcheck = &healthConfig
	}
}

func WithImageHealthcheck(healthConfig dockerContainer.HealthConfig) MockContainerUpdate {
	return func(cnt dockerContainer.InspectResponse, img image.InspectResponse) {
		img.Config.Healthcheck = &healthConfig
	}
}
