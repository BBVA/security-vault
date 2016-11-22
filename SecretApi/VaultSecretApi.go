package SecretApi

import (
	"bytes"
	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/rancher/secrets-bridge/pkg/archive"
	"path/filepath"
	"time"
)

type VaultSecretApi struct {
	client             *vault.Client
	role               string
	persistenceChannel chan persistence.LeaseEvent
	config             config.ConfigHandler
}

func NewVaultSecretApi(mainConfig config.ConfigHandler, persistenceChannel chan persistence.LeaseEvent) (*VaultSecretApi, error) {

	vaultConfig := vault.DefaultConfig()

	if err := vaultConfig.ReadEnvironment(); err != nil {
		return nil, err
	}

	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	token, err := mainConfig.GetToken()
	if err != nil {
		return nil, err
	}

	client.SetToken(token)
	client.SetAddress(mainConfig.GetVaultServer())

	return &VaultSecretApi{
		client:             client,
		role:               mainConfig.GetRole(),
		persistenceChannel: persistenceChannel,
		config:             mainConfig,
	}, nil
}

func (api *VaultSecretApi) GetSecretFiles(commonName string, containerID string) (*bytes.Buffer, error) {
	fmt.Println("Generating secret\n")
	files := []archive.ArchiveFile{}
	params := make(map[string]interface{})
	params["common_name"] = commonName

	path := filepath.Join("pki/issue/", api.role)

	secrets, err := api.client.Logical().Write(path, params)
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

	timestamp := time.Now().Unix()

	api.persistenceChannel <- persistence.LeaseEvent{
		EventType:   "start",
		ContainerID: containerID,
		Lease: persistence.LeaseInfo{
			LeaseID:   secrets.LeaseID,
			LeaseTime: secrets.LeaseDuration,
			Renewable: secrets.Renewable,
			Timestamp: timestamp,
		},
	}

	return tarball, nil

}

func (api *VaultSecretApi) DeleteSecrets(containerID string) error {
	fmt.Println("Deleting secret persistence..")
	event := persistence.LeaseEvent{
		EventType:   "stop",
		ContainerID: containerID,
		Lease:       persistence.LeaseInfo{},
	}
	api.persistenceChannel <- event

	return nil
}
