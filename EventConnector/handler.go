package EventConnector

import (
	"bytes"
	"fmt"
	"time"

	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"descinet.bbva.es/cloudframe-security-vault/utils/archive"
	. "github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"golang.org/x/net/context"
)

func (c *DockerConnector) eventHandler(msg *events.Message) {
	fmt.Printf("Received:\n Action %v\nActor %v\nFrom %v\nID %v\nStatus %v\nTime %v\nTimenano %v\nType %v\n", msg.Action, msg.Actor, msg.From, msg.ID, msg.Status, msg.Time, msg.TimeNano, msg.Type)
	switch msg.Action {
	case "start":
		id, ok := msg.Actor.Attributes["common_name"]
		if ok {
			fmt.Println("label detected!")
			if err := c.CopySecretsToContainer(id, msg.ID); err != nil {
				panic(err.Error())
			}

		}
	case "stop":
		if err := c.secretApiHandler.DeleteSecrets(msg.ID); err != nil {
			panic(err.Error())
		}

		event := persistence.LeaseEvent{
			EventType:  "stop",
			Identifier: msg.ID,
			Lease:      persistence.LeaseInfo{},
		}
		c.persistenceChannel <- event
	}
}

func (c *DockerConnector) CopySecretsToContainer(common_name string,containerID string ) error {

	secrets, err := c.secretApiHandler.GetSecretFiles(common_name)
	if err != nil {
		return err
	}

	tarball, err := secretsToTarball(secrets)
	if err != nil {
		return err
	}

	opts := CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	}
	if err := c.cli.CopyToContainer(context.Background(), containerID, c.path, tarball, opts); err != nil {
		return err
	}

	timestamp := time.Now().Unix()

	c.persistenceChannel <- persistence.LeaseEvent{
		EventType:  "start",
		Identifier: containerID,
		Lease: persistence.LeaseInfo{
			CommonName: common_name,
			LeaseID:   secrets.LeaseID,
			LeaseTime: secrets.LeaseDuration,
			Renewable: secrets.Renewable,
			Timestamp: timestamp,
		},
	}
	return nil
}

func secretsToTarball(secrets *SecretApi.Secrets) (*bytes.Buffer, error) {
	files := []archive.ArchiveFile{}
	files = append(files, archive.ArchiveFile{Name: "private", Content: secrets.Private})
	files = append(files, archive.ArchiveFile{Name: "cacert", Content: secrets.Cacert})
	files = append(files, archive.ArchiveFile{Name: "public", Content: secrets.Public})

	tarball, err := archive.CreateTarArchive(files)
	if err != nil {
		return nil, err
	}

	return tarball, nil
}
