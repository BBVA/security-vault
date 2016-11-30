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
	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"io"
)

type Connector interface {
	StartConnector() error
	eventHandler(msg *events.Message)
}

type DockerClient interface {
	Events(ctx context.Context, options types.EventsOptions) (io.ReadCloser, error)
	CopyToContainer(ctx context.Context, container, path string, content io.Reader, options types.CopyToContainerOptions) error
}

type DockerConnector struct {
	secretApiHandler SecretApi.SecretApi
	cli              DockerClient
	path             string
	dockerClient     func() (DockerClient, error)
	persistenceChannel chan persistence.LeaseEvent
}

func NewConnector(secretApiHandler SecretApi.SecretApi, config config.ConfigHandler, persistenceChannel chan persistence.LeaseEvent) *DockerConnector {
	return &DockerConnector{
		secretApiHandler: secretApiHandler,
		path:             config.GetSecretPath(),
		dockerClient:     getDockerClient,
		persistenceChannel: persistenceChannel,
	}
}

func (c *DockerConnector) Start() error {

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

func getDockerClient() (DockerClient, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	return client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
}
