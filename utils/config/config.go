package config

import (
	"fmt"
	"errors"
	"os"
)

type Config map[string]string

func ReadConfig() (Config, error) {

	fmt.Println("Reading environment...\n")

	configMap := make(map[string]string)

	configMap["vaultServer"] = os.Getenv("VAULT_SERVER")
	configMap["tokenPath"] = os.Getenv("TOKEN_PATH")
	configMap["secretPath"] = os.Getenv("SECRET_PATH")
	configMap["role"] = os.Getenv("ROLE")
	//configMap["persistencePath"] = os.Getenv("PERSISTENCE_PATH")

	for k, v := range configMap {
		if len(v) == 0 {
			err := fmt.Sprintf("Undefined configuration: %s",k)
			return nil, errors.New(err)
		}
	}
return configMap, nil
}

