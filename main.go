package main

import (
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	"github.com/facebookgo/inject"
	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
)

func main() {

	cfg := &config.Config{}
	if err := inject.Populate(cfg, &filesystem.DefaultFileUtils{}); err != nil {
		panic(err.Error())
	}

	if err := cfg.ReadConfig(); err != nil {
		panic(err.Error())
	}

	secretApiHandler, err := SecretApi.NewVaultSecretApi(cfg)
	if err != nil {
		panic(err.Error())
	}

	persistenceChannel, persistenceManager := persistence.NewPersistenceManager(cfg)
	if err := persistenceManager.RecoverLeases(); err != nil {
		panic(err.Error())
	}

	go persistenceManager.Run()

	connector := EventConnector.NewConnector(secretApiHandler, cfg, persistenceChannel)

	connector.Start()

}
