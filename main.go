package main

import (
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"github.com/facebookgo/inject"
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

	persistenceCfg := &persistence.PersistenceManager{}
	if err := inject.Populate(persistenceCfg, &filesystem.DefaultFileUtils{}); err != nil {
		panic(err.Error())
	}
	persistenceChannel, persistenceManager := persistence.NewPersistenceManager(cfg, persistenceCfg)
	if err := persistenceManager.RecoverLeases(); err != nil {
		panic(err.Error())
	}

	go persistenceManager.Run()

	cli, err := EventConnector.GetDockerClient()
	if err != nil {
		panic(err.Error())
	}

	connector := EventConnector.NewConnector(secretApiHandler, cfg, cli, persistenceChannel)

	connector.Start()

}
