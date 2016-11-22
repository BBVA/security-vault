package main

import (
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
)

func main() {

	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err.Error())
	}

	persistenceChannel, persistenceManager := persistence.NewPersistenceManager(cfg)

	secretApiHandler, err := SecretApi.NewVaultSecretApi(cfg)
	if err != nil {
		panic(err.Error())
	}

	go persistenceManager.Run()

	connector := EventConnector.NewConnector(secretApiHandler, cfg, persistenceChannel)

	connector.StartConnector()

}
