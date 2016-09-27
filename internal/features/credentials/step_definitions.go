package credentials

import (
	. "github.com/gucumber/gucumber"
	"github.com/docker/docker/client"
	"context"
	"github.com/docker/docker/api/types/container"
	"fmt"
	"strings"
	"github.com/docker/docker/api/types"
)

type Container struct {
	name    string
	options []string
	id      string
}

func init() {

	containers := make(map[string]*Container)

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}

	After("@destroyContainer", func() {
		for _, c := range containers {
			fmt.Println(c.id)
			if err := cli.ContainerRemove(context.Background(), c.id, types.ContainerRemoveOptions{
				RemoveLinks: false,
				RemoveVolumes: true,
				Force: true,
			}); err != nil {
				panic(err)
			}

		}
	})

	Given(`^a container "(.+?)" configured with the following volume driver options:$`, func(containerName string, volDriverOpts [][]string) {
		volumeOption := volDriverOpts[1]
		containers[containerName] = &Container{containerName, volumeOption, ""}
	})

	When(`^the container "(.+?)" is started$`, func(containerName string) {
		c := containers[containerName]

		vols := make(map[string]struct{})
		vols[c.options[2]] = struct{}{}

		containerConfig := &container.Config{
			Cmd: strings.Split("cat " + c.options[2] + "/credential", " "),
			Image: "alpine",
			Volumes: vols,

		}

		hostConfig := &container.HostConfig{
			AutoRemove: true,
			VolumeDriver: c.options[0],
			// TODO Mounts:
		}

		fmt.Println(containerConfig)
		fmt.Println(hostConfig)

		response, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, containerName)
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(context.Background(), response.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		containers[containerName].id = response.ID

		fmt.Println(response.ID)

	})

	Then(`^the container "(.+?)" credentials will be the following$`, func(containerName string, credentialInfo [][]string) {
		T.Skip() // pending
	})

}
