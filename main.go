package main

import (
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
)


func main() {

	config,err := config.ReadConfig()
	if err != nil {
		panic(err.Error())
	}

	secretApiHandler, err := SecretApi.NewVaultSecretApi(config)
	if err != nil {
		panic(err.Error())
	}

	go secretApiHandler.PersistenceManager()

	connector := EventConnector.NewConnector(secretApiHandler,config)
	connector.StartConnector()

}
