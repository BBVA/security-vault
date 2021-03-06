package credentials

import (
	"context"
	"strings"

	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	. "github.com/gucumber/gucumber"
	"time"
)

type Container struct {
	name string
	id   string
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

	Given(`^a container "(.+?)" configured$`, func(containerName string) {

		containers[containerName] = &Container{containerName, ""}
	})

	When(`^the container "(.+?)" is started$`, func(containerName string) {

		containerConfig := createContainerConfiguration()

		hostConfig := createHostConfiguration()

		containers[containerName].id = runContainer(cli, containerName, hostConfig, containerConfig)

		time.Sleep(10 * time.Second)

	})

	Then(`^the container "(.+?)" credentials will be the following$`, func(containerName string, credentialInfo [][]string) {

		for _, cred := range credentialInfo[1:] {
			expectedContent := strings.TrimSpace(cred[1])
			fmt.Printf("nombre del fichero: %s\n", cred[0])
			content := strings.TrimSpace(getFileContent(cli, containerName, cred[0]))

			//fmt.Println(expectedContent, len(expectedContent))
			//fmt.Println(content, len(content))

			//if !reflect.DeepEqual(expectedContent, content)
			if !strings.Contains(content, expectedContent) {
				T.Errorf("Expected: " + expectedContent + " Actual: " + content)
			}
		}
	})
}

func createDockerClient() *client.Client {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	if err != nil {
		T.Errorf(err.Error())
	}

	return cli
}

func createContainerConfiguration() *container.Config {

	cmd := []string{"/bin/ash", "-c", "while true; do sleep 5; done"}

	labels := make(map[string]string)
	labels["common_name"] = "makecloudframegreatagain.cloudframe.wtf"

	return &container.Config{
		Cmd:    cmd,
		Image:  "alpine",
		Labels: labels,
	}
}

func createHostConfiguration() *container.HostConfig {
	return &container.HostConfig{
		AutoRemove: true,
	}
}

func runContainer(cli *client.Client, containerName string, hostConfiguration *container.HostConfig, containerConfiguration *container.Config) string {
	response, err := cli.ContainerCreate(context.Background(), containerConfiguration, hostConfiguration, nil, containerName)
	if err != nil {
		T.Errorf(err.Error())
	}

	if err := cli.ContainerStart(context.Background(), response.ID, types.ContainerStartOptions{}); err != nil {
		T.Errorf(err.Error())
	}

	return response.ID
}

func getFileContent(cli *client.Client, container string, fileName string) string {

	cmd := []string{"/bin/cat", fileName}

	options := types.ExecConfig{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}
	execresponse, err := cli.ContainerExecCreate(context.Background(), container, options)
	if err != nil {
		T.Errorf(err.Error())
	}

	connection, err := cli.ContainerExecAttach(context.Background(), execresponse.ID, options)
	if err != nil {
		T.Errorf(err.Error())
	}

	defer connection.Close()

	output, err := connection.Reader.ReadString('\n')
	if err != nil {
		T.Errorf(err.Error())
	}
	if len(output) == 0 {
		T.Errorf("no data returned\n")
	}

	return output[8:]
}

func destroyContainer(cli *client.Client, containers map[string]*Container) {
	for _, c := range containers {
		if err := cli.ContainerRemove(context.Background(), c.id, types.ContainerRemoveOptions{
			RemoveLinks:   false,
			RemoveVolumes: true,
			Force:         true,
		}); err != nil {
			T.Errorf(err.Error())
		}
	}
}
