package main

import (
	"syscall"
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
	"golang.org/x/sys/unix"
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

	if err := syscall.Mlockall(unix.MCL_FUTURE | unix.MCL_CURRENT); err != nil {
		panic(err.Error())
	}

	exampleSecretApiHandler, err := SecretApi.NewExampleSecretApi(DefaultCaCert, DefaultPrivateKey, DefaultPublicKey, &fileUtils)
	if err != nil {
		panic(err.Error())
	}

	connector := EventConnector.NewConnector(exampleSecretApiHandler, DefaultSecretPaths)
	connector.StartConnector()

}
