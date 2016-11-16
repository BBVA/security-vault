package SecretApi

import (
	"bytes"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/rancher/secrets-bridge/pkg/archive"
	"io/ioutil"
	"path/filepath"
)

type leaseInfo struct {
	leaseID   string
	leaseTime int
	renewable bool
}
type leaseEvent struct {
	eventType   string
	containerID string
	lease       leaseInfo
}

type VaultSecretApi struct {
	client             *vault.Client
	role               string
	leases             map[string]leaseInfo
	persistenceChannel chan leaseEvent
}

func NewVaultSecretApi(mainConfig config.Config) (*VaultSecretApi, error) {

	config := vault.DefaultConfig()

	if err := config.ReadEnvironment(); err != nil {
		return nil, err
	}

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	token, err := ioutil.ReadFile(mainConfig["tokenPath"])
	if err != nil {
		return nil, err
	}

	client.SetToken(string(token))
	client.SetAddress(mainConfig["vaultServer"])

	leases := make(map[string]leaseInfo)
	persistenceChannel := make(chan leaseEvent)

	return &VaultSecretApi{
		client:             client,
		role:               mainConfig["role"],
		leases:             leases,
		persistenceChannel: persistenceChannel,
	}, nil
}

func (Api *VaultSecretApi) GetSecretFiles(commonName string, containerID string) (*bytes.Buffer, error) {

	files := []archive.ArchiveFile{}
	params := make(map[string]interface{})
	params["common_name"] = commonName

	path := filepath.Join("pki/issue/", Api.role)

	secrets, err := Api.client.Logical().Write(path, params)
	if err != nil {
		return nil, err
	}

	files = append(files, archive.ArchiveFile{Name: "private", Content: secrets.Data["private_key"].(string)})
	files = append(files, archive.ArchiveFile{Name: "cacert", Content: secrets.Data["issuing_ca"].(string)})
	files = append(files, archive.ArchiveFile{Name: "public", Content: secrets.Data["certificate"].(string)})

	tarball, err := archive.CreateTarArchive(files)
	if err != nil {
		return nil, err
	}

	Api.persistenceChannel <- leaseEvent{
		eventType:   "start",
		containerID: containerID,
		lease: leaseInfo{
			leaseID:   secrets.LeaseID,
			leaseTime: secrets.LeaseDuration,
			renewable: secrets.Renewable,
		},
	}

	return tarball, nil

}

func (Api *VaultSecretApi) DeleteSecrets(containerID string) error{
	event := leaseEvent{
		eventType: "stop",
		containerID: containerID,
		lease: leaseInfo{},
	}
	Api.persistenceChannel <- event

	return nil
}

func (Api *VaultSecretApi) PersistenceManager() {
	select {
	case <-Api.persistenceChannel:
		fmt.Println("Lease event received\n")
		var event leaseEvent

		event = <-Api.persistenceChannel
		switch event.eventType {
		case "start":
			Api.leases[event.containerID] = event.lease
		case "stop":
			_, ok := Api.leases[event.containerID]
			if ok {
				delete(Api.leases,event.containerID)
			}
		}

	}

}
