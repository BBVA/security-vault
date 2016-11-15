package SecretApi

import (
	vault "github.com/hashicorp/vault/api"
	"bytes"
	"github.com/rancher/secrets-bridge/pkg/archive"
	"path/filepath"
	"io/ioutil"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
)

type VaultSecretApi struct {
	client *vault.Client
	role string
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

	return &VaultSecretApi{
		client: client,
		role: mainConfig["role"],
	},nil
}

func (Api *VaultSecretApi) GetSecretFiles(commonName string) (*bytes.Buffer, error) {

	files := []archive.ArchiveFile{}
	params := make(map[string]interface{})
	params["common_name"] = commonName

	path := filepath.Join("pki/issue/", Api.role)

	secrets, err := Api.client.Logical().Write(path, params)
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