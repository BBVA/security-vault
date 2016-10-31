package main

import (
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"path/filepath"
	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
)

var (
	DefaultPublicKey  = filepath.Join("/tmp", "public.key")
	DefaultPrivateKey = filepath.Join("/tmp", "private.key")
	DefaultCaCert     = filepath.Join("/tmp", "ca.cert")
	DefaultSecretPaths = "/tmp"
)


func main() {

	fileUtils := filesystem.DefaultFileUtils{}


	exampleSecretApiHandler, err := SecretApi.NewExampleSecretApi(DefaultCaCert, DefaultPrivateKey, DefaultPublicKey, &fileUtils)
	if err != nil {
		panic(err.Error())
	}

	connector := EventConnector.NewConnector(exampleSecretApiHandler, DefaultSecretPaths)
	connector.StartConnector()

}
