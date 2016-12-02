package EventConnector

import (
	"encoding/json"
	"log"

	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"

	"io"

	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
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
	secretApiHandler   SecretApi.SecretApi
	cli                DockerClient
	path               string
	eventStream        io.ReadCloser
	persistenceChannel chan persistence.LeaseEvent
}

func NewConnector(secretApiHandler SecretApi.SecretApi, config config.ConfigHandler, client DockerClient, persistenceChannel chan persistence.LeaseEvent) (*DockerConnector, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("event", "start")
	filterArgs.Add("event", "stop")

	eventOptions := types.EventsOptions{
		Filters: filterArgs,
	}

	eventsResp, err := client.Events(context.Background(), eventOptions)
	if err != nil {
		return nil, err
	}

	return &DockerConnector{
		secretApiHandler:   secretApiHandler,
		cli:                client,
		path:               config.GetSecretPath(),
		eventStream:        eventsResp,
		persistenceChannel: persistenceChannel,
	}, nil
}

func (c *DockerConnector) Start() error {
	log.Println("Entering event listening Loop")
	d := json.NewDecoder(c.eventStream)
	for {
		msg := &events.Message{}
		err := d.Decode(msg)
		if err != nil {
			return err
		}

		go c.eventHandler(msg)

	}
}

func (c *DockerConnector) Stop() error {
	return c.eventStream.Close()
}

func GetDockerClient() (DockerClient, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	return client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
}
