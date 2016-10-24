package main

import (
	"fmt"
	"path/filepath"

	"syscall"

	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"descinet.bbva.es/cloudframe-security-vault/utils/fuse"
	"github.com/docker/go-plugins-helpers/volume"
	"golang.org/x/sys/unix"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
)

const (
	SocketAddress = "/run/docker/plugins/Vault.sock"
	ServerUrl     = "http://localhost:8200"
	VaultToken    = ""
)

var (
	//DefaultPath = filepath.Join(volume.DefaultDockerRootDirectory, "_vault")
	DefaultMountPath = filepath.Join("/mnt/volumes", "_vault")
	DefaultConfigPath = filepath.Join("/tmp", "security-vault")
)

/* hay que implementar argumentos para recibir:
*
*  -path de la configuración ( aquí pondremos todo lo que ahora está como constantes )
*
 */

func main() {

	if err := syscall.Mlockall(unix.MCL_FUTURE | unix.MCL_CURRENT); err != nil {
		panic(err.Error())
	}

	fuse := filesystem.DefaultFuseWrapper{}
	dirUtils := filesystem.DefaultDirUtils{}
	fileUtils := filesystem.DefaultFileUtils{}
	ExampleSecretApiHandler := SecretApi.NewExampleSecretApi()
	secretApiHandler := SecretApi.NewSecretApi(ExampleSecretApiHandler)
	fuseUtils := fuseutils.NewFuseUtils(fuse, secretApiHandler)
	driver := NewVaultDriver(DefaultMountPath, ServerUrl, VaultToken, &dirUtils, fuseUtils)
	persitor, _ := NewVolumePersistor(DefaultConfigPath, driver, &dirUtils, &fileUtils)
	handler := volume.NewHandler(persitor)

	fmt.Printf("Listening on %s\n", SocketAddress)
	fmt.Println(handler.ServeUnix("root", SocketAddress))
}

