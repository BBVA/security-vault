package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigHandler interface {
	GetToken() (string, error)
	GetRole() string
	GetVaultServer() string
	GetPersistencePath() string
	GetSecretPath() string
	Get(string) (string, error)
}

type Config map[string]string

func ReadConfig() (Config, error) {

	fmt.Println("Reading environment...\n")

	configMap := make(map[string]string)

	configMap["vaultServer"] = os.Getenv("VAULT_SERVER")
	configMap["tokenPath"] = os.Getenv("TOKEN_PATH")
	configMap["secretPath"] = os.Getenv("SECRET_PATH")
	configMap["role"] = os.Getenv("ROLE")
	configMap["persistencePath"] = os.Getenv("PERSISTENCE_PATH")
	configMap["socket"] = os.Getenv("SOCKET")

	return validateConfiguration(configMap)
}

func (c Config) GetToken() (string, error) {
	if token, err := ioutil.ReadFile(c["tokenPath"]); err != nil {
		return "", err
	} else {

		return string(token), nil
	}
}

func (c Config) GetRole() string {
	return c["role"]
}

func (c Config) GetVaultServer() string {
	return c["vaultServer"]
}

func (c Config) GetPersistencePath() string {
	return c["persistencePath"]
}

func (c Config) GetSecretPath() string {
	return c["secretPath"]
}

func (c Config) Get(key string) (string, error) {
	if value, ok := c[key]; ok {
		return value, nil
	} else {
		return "", errors.New(fmt.Sprintf("Missing Key: %s", key))
	}
}

func validateConfiguration(cfg Config) (Config, error) {
	for k, v := range cfg {
		if len(v) == 0 {
			err := fmt.Sprintf("Undefined configuration: %s", k)
			return nil, errors.New(err)
		}
	}

	return cfg, nil
}
