package EventConnector

import (
	. "github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"golang.org/x/net/context"

	"log"
)

func (c *DockerConnector) eventHandler(msg *events.Message) {
	log.Printf("Received:\n Action %v\nActor %v\nFrom %v\nID %v\nStatus %v\nTime %v\nTimenano %v\nType %v\n", msg.Action, msg.Actor, msg.From, msg.ID, msg.Status, msg.Time, msg.TimeNano, msg.Type)

	id, ok := msg.Actor.Attributes["common_name"]
	if ok {
		log.Println("label detected!")
		tarball,err := c.secretApiHandler.GetSecretFiles(id)
		if err != nil {
			log.Printf("GetSecretFiles failed %s\n", err.Error())
		}
		opts := CopyToContainerOptions{
			AllowOverwriteDirWithFile: false,
		}
		if err := c.cli.CopyToContainer(context.Background(), msg.ID, c.path, tarball, opts); err != nil {
			log.Printf("CopyToContainer failed: %s\n", err.Error())
		}
	}
}
