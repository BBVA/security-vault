package SecretApi

import (
	"fmt"
	"path/filepath"

	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	vault "github.com/hashicorp/vault/api"
)

type VaultSecretApi struct {
	client *vault.Client
	role   string
	config config.ConfigHandler
}

func NewVaultSecretApi(mainConfig config.ConfigHandler) (*VaultSecretApi, error) {

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
		client: client,
		role:   mainConfig.GetRole(),
		config: mainConfig,
	}, nil
}

func (api *VaultSecretApi) GetSecretFiles(commonName string) (*Secrets, error) {
	fmt.Println("Generating secret")
	params := make(map[string]interface{})
	params["common_name"] = commonName

	path := filepath.Join("pki/issue/", api.role)

	secrets, err := api.client.Logical().Write(path, params)
	if err != nil {
		return nil, err
	}

	return &Secrets{
		Public:        secrets.Data["certificate"].(string),
		Private:       secrets.Data["private_key"].(string),
		Cacert:        secrets.Data["issuing_ca"].(string),
		LeaseID:       secrets.LeaseID,
		LeaseDuration: secrets.LeaseDuration,
		Renewable:     secrets.Renewable,
	}, nil

}

func (api *VaultSecretApi) DeleteSecrets(containerID string) error {
	fmt.Println("Deleting secret persistence..")
	return nil
}
