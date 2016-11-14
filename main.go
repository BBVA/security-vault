package main

import (
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
)

var (
	DefaultSecretPaths = "/tmp"
)


func main() {

	secretApiHandler, err := SecretApi.NewVaultSecretApi("86666040-1b49-35e8-5bb7-e4c323f48df3","http://vault-server:8200")
	if err != nil {
		panic(err.Error())
	}

	connector := EventConnector.NewConnector(secretApiHandler, DefaultSecretPaths)
	connector.StartConnector()

}
