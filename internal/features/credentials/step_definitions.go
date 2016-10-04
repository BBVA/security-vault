package credentials

import (
	"bytes"
	"context"
	"errors"
	//"fmt"
	"reflect"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	. "github.com/gucumber/gucumber"
)

type Container struct {
	name    string
	options []string
	id      string
}

func init() {

	containers := make(map[string]*Container)
	cli := createDockerClient()

	BeforeAll(func() {
	})

	AfterAll(func() {
	})

	After("@destroyContainers", func() {
		destroyContainer(cli, containers)
	})

	Given(`^a container "(.+?)" configured with the following volume driver options:$`, func(containerName string, volDriverOpts [][]string) {
		volumeOption := volDriverOpts[1]
		containers[containerName] = &Container{containerName, volumeOption, ""}
	})

	When(`^the container "(.+?)" is started$`, func(containerName string) {
		c := containers[containerName]
		volumeDriver := c.options[0]
		hostFS := c.options[1]
		containerMountPoint := c.options[2]

		vols := make(map[string]struct{})
		vols[c.options[2]] = struct{}{}

		containerConfig := createContainerConfiguration(containerMountPoint)

		hostConfig := createHostConfiguration(volumeDriver, hostFS, containerMountPoint)

		containers[containerName].id = runContainer(cli, containerName, hostConfig, containerConfig)

	})

	Then(`^the container "(.+?)" credentials will be the following$`, func(containerName string, credentialInfo [][]string) {
		containerId := containers[containerName].id

		expectedContent := strings.TrimSpace(credentialInfo[1][1])
		content := strings.TrimSpace(getContainerLogs(cli, containerId))

		//fmt.Println(expectedContent, len(expectedContent))
		//fmt.Println(content, len(content))

		if !reflect.DeepEqual(expectedContent, content) {
			panic(errors.New("Expected: " + expectedContent + " Actual: " + content))
		}
	})
}

func createDockerClient() *client.Client {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}

	return cli
}

func createContainerConfiguration(volume string) *container.Config {
	vols := make(map[string]struct{})
	vols[volume] = struct{}{}

	return &container.Config{
		Cmd:     strings.Split("cat "+volume+"/credential", " "),
		Image:   "alpine",
		Volumes: vols,
	}
}

func createHostConfiguration(volumeDriver, hostFS, containerMountPoint string) *container.HostConfig {
	return &container.HostConfig{
		AutoRemove:   true,
		VolumeDriver: volumeDriver,
		Binds: []string{
			strings.Join([]string{hostFS, containerMountPoint}, ":"),
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Target: hostFS,
				BindOptions: &mount.BindOptions{
					Propagation: mount.PropagationRPrivate,
				},
				VolumeOptions: &mount.VolumeOptions{
					DriverConfig: &mount.Driver{
						Name: hostFS,
					},
				},
			},
		},
	}
}

func runContainer(cli *client.Client, containerName string, hostConfiguration *container.HostConfig, containerConfiguration *container.Config) string {
	response, err := cli.ContainerCreate(context.Background(), containerConfiguration, hostConfiguration, nil, containerName)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(context.Background(), response.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return response.ID
}

func getContainerLogs(cli *client.Client, containerId string) string {
	reader, err := cli.ContainerLogs(context.Background(), containerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: false,
		Timestamps: false,
	})

	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	//l, err := buf.ReadFrom(reader)
	buf.ReadFrom(reader)
	//fmt.Printf("read %d bytes\n", l)
	return buf.String()[8:]
}

func destroyContainer(cli *client.Client, containers map[string]*Container) {
	for _, c := range containers {
		if err := cli.ContainerRemove(context.Background(), c.id, types.ContainerRemoveOptions{
			RemoveLinks:   false,
			RemoveVolumes: true,
			Force:         true,
		}); err != nil {
			panic(err)
		}
	}
}
