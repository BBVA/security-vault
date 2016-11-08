package EventConnector

import (
	"encoding/json"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
	"log"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"

)

type Connector interface {
	StartConnector() error
	eventHandler(msg *events.Message)
}

type DockerConnector struct{
	secretApiHandler SecretApi.SecretApi
	cli *client.Client
	path string
}

func NewConnector(secretApiHandler SecretApi.SecretApi, path string) *DockerConnector {
	return &DockerConnector{
		secretApiHandler: secretApiHandler,
		path: path,
	}
}

func (c *DockerConnector) StartConnector() error {

	cli, err := getDockerClient()
	if err != nil {
		log.Printf("Could not get Docker client: %s", err)
	}
	c.cli = cli

	filterArgs := filters.NewArgs()
	filterArgs.Add("event", "start")

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
