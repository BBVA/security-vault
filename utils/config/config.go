package config

import (
	"errors"
	"fmt"
	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
)

type ConfigHandler interface {
	ReadConfig() error
	GetToken() (string, error)
	GetRole() string
	GetVaultServer() string
	GetPersistencePath() string
	GetSecretPath() string
	Get(string) (string, error)
}

type Config struct {
	cfg       map[string]string
	FileUtils filesystem.FileUtils `inject:""`
}

func (c *Config) ReadConfig() error {

	fmt.Println("Reading environment...")

	configMap := make(map[string]string)
	
	configMap["vaultServer"] = c.FileUtils.Getenv("VAULT_SERVER")
	configMap["tokenPath"] = c.FileUtils.Getenv("TOKEN_PATH")
	configMap["secretPath"] = c.FileUtils.Getenv("SECRET_PATH")
	configMap["role"] = c.FileUtils.Getenv("ROLE")
	configMap["persistencePath"] = c.FileUtils.Getenv("PERSISTENCE_PATH")

	c.cfg = configMap

	return c.validateConfiguration()
}

func (c Config) GetToken() (string, error) {
	value, _ := c.Get("tokenPath")
	if token, err := c.FileUtils.ReadFile(value); err != nil {
		return "", err
	} else {

		return string(token), nil
	}
}

func (c Config) GetRole() string {
	value, _ := c.Get("role")
	return value
}

func (c Config) GetVaultServer() string {
	value, _ := c.Get("vaultServer")
	return value
}

func (c Config) GetPersistencePath() string {
	value, _ := c.Get("persistencePath")
	return value
}

func (c Config) GetSecretPath() string {
	value, _ := c.Get("secretPath")
	return value
}

func (c Config) Get(key string) (string, error) {
	if value, ok := c.cfg[key]; ok {
		return value, nil
	} else {
		return "", errors.New(fmt.Sprintf("Missing Key: %s", key))
	}
}

func (c Config)validateConfiguration() error {
	for k, v := range c.cfg {
		if len(v) == 0 {
			return errors.New(fmt.Sprintf("Undefined configuration: %s", k))
		}
	}

	return nil
}
