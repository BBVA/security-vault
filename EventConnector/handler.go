package EventConnector

import (
	"fmt"
	. "github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"golang.org/x/net/context"

	"log"
)

func (c *DockerConnector) eventHandler(msg *events.Message) {
	log.Printf("Received:\n Action %v\nActor %v\nFrom %v\nID %v\nStatus %v\nTime %v\nTimenano %v\nType %v\n", msg.Action, msg.Actor, msg.From, msg.ID, msg.Status, msg.Time, msg.TimeNano, msg.Type)
	if msg.Actor.Attributes["name"] == "cred-test" {
		tarball := c.secretApiHandler.GetSecretFiles()

		opts := CopyToContainerOptions{
			AllowOverwriteDirWithFile: false,
		}
		if err := c.cli.CopyToContainer(context.Background(), msg.ID, c.path, tarball, opts); err != nil {
			fmt.Println("CopyToContainer failed: %s", err.Error())
		}

	}
}
