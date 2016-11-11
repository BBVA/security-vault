package SecretApi

import (
	vault "github.com/hashicorp/vault/api"
	"bytes"
	"github.com/rancher/secrets-bridge/pkg/archive"
)

type VaultSecretApi struct {
	client *vault.Client
}

func NewVaultSecretApi(token string, url string) (*VaultSecretApi, error) {

	config := vault.DefaultConfig()

	if err := config.ReadEnvironment(); err != nil {
		return nil, err
	}

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)
	client.SetAddress(url)

	return &VaultSecretApi{
		client: client,
	},nil
}

func (Api *VaultSecretApi) GetSecretFiles(commonName string) (*bytes.Buffer, error) {

	files := []archive.ArchiveFile{}
	params := make(map[string]interface{})
	params["common_name"] = commonName

	secrets, err := Api.client.Logical().Write("pki/issue/cloudframe-dot-wtf", params)
	if err != nil {
		return nil, err
	}

	files = append(files,archive.ArchiveFile{ Name: "private",Content: secrets.Data["private_key"].(string)})
	files = append(files,archive.ArchiveFile{ Name: "cacert",Content: secrets.Data["issuing_ca"].(string)})
	files = append(files,archive.ArchiveFile{ Name: "public",Content: secrets.Data["certificate"].(string)})

	tarball, err := archive.CreateTarArchive(files)
	if err != nil {
		return nil,err
	}
	return tarball, nil

}