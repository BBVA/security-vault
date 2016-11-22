package EventConnector

import (
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"encoding/json"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
	"log"

	"descinet.bbva.es/cloudframe-security-vault/utils/config"
)

type Connector interface {
	StartConnector() error
	eventHandler(msg *events.Message)
}

type DockerConnector struct {
	secretApiHandler SecretApi.SecretApi
	cli              *client.Client
	path             string
	dockerClient     func() (*client.Client, error)
}

func NewConnector(secretApiHandler SecretApi.SecretApi, config config.Config) *DockerConnector {
	return &DockerConnector{
		secretApiHandler: secretApiHandler,
		path:             config["secretPath"],
		dockerClient:     getDockerClient,
	}
}

func (c *DockerConnector) StartConnector() error {

	cli, err := c.dockerClient()
	if err != nil {
		log.Printf("Could not get Docker client: %s", err)
	}
	c.cli = cli

	filterArgs := filters.NewArgs()
	filterArgs.Add("event", "start")
	filterArgs.Add("event", "stop")

	eventOptions := types.EventsOptions{
		Filters: filterArgs,
	}

	eventsResp, err := cli.Events(context.Background(), eventOptions)
	if err != nil {
		return err
	}
	defer eventsResp.Close()

	log.Println("Entering event listening Loop")
	d := json.NewDecoder(eventsResp)
	for {
		msg := &events.Message{}
		d.Decode(msg)

		go c.eventHandler(msg)

	}

}

func getDockerClient() (*client.Client, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	return client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
}
