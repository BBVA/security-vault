package main

import (
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"path/filepath"
	//"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
)

var (
	DefaultPublicKey  = filepath.Join("/tmp", "public.key")
	DefaultPrivateKey = filepath.Join("/tmp", "private.key")
	DefaultCaCert     = filepath.Join("/tmp", "ca.cert")
	DefaultSecretPaths = "/tmp"
)


func main() {

	//fileUtils := filesystem.DefaultFileUtils{}


	//secretApiHandler, err := SecretApi.NewExampleSecretApi(DefaultCaCert, DefaultPrivateKey, DefaultPublicKey, &fileUtils)
	secretApiHandler, err := SecretApi.NewVaultSecretApi("86666040-1b49-35e8-5bb7-e4c323f48df3","http://vault-server:8200")
	if err != nil {
		panic(err.Error())
	}

	connector := EventConnector.NewConnector(secretApiHandler, DefaultSecretPaths)
	connector.StartConnector()

}
