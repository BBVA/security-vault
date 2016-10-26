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
	"log"
)

const (
	SocketAddress = "/run/docker/plugins/Vault.sock"
	ServerUrl = "http://localhost:8200"
	VaultToken = ""
)

var (
	//DefaultPath = filepath.Join(volume.DefaultDockerRootDirectory, "_vault")
	DefaultMountPath  = filepath.Join("/mnt/volumes", "_vault")
	DefaultConfigPath = filepath.Join("/tmp", "security-vault")
	DefaultPublicKey  = filepath.Join("/tmp", "public.key")
	DefaultPrivateKey = filepath.Join("/tmp", "private.key")
	DefaultCaCert     = filepath.Join("/tmp", "ca.cert")
)

/* hay que implementar argumentos para recibir:
*
*  -path de la configuración ( aquí pondremos todo lo que ahora está como constantes )
*
 */

func main() {

	fuse := fuseutils.DefaultFuseWrapper{}
	dirUtils := filesystem.DefaultDirUtils{}
	fileUtils := filesystem.DefaultFileUtils{}

	if err := syscall.Mlockall(unix.MCL_FUTURE | unix.MCL_CURRENT); err != nil {
		panic(err.Error())
	}

	exampleSecretApiHandler, err := SecretApi.NewExampleSecretApi(DefaultCaCert, DefaultPrivateKey, DefaultPublicKey, &fileUtils)
	if err != nil {
		panic(err.Error())
	}

	fuseUtils := fuseutils.NewFuseUtils(fuse, exampleSecretApiHandler)
	driver := NewVaultDriver(DefaultMountPath, ServerUrl, VaultToken, &dirUtils, fuseUtils)
	persitor, err := NewVolumePersistor(DefaultConfigPath, &driver, &dirUtils, &fileUtils)
	if err != nil {
		log.Panic(err)
	}
	handler := volume.NewHandler(persitor)

	fmt.Printf("Listening on %s\n", SocketAddress)
	fmt.Println(handler.ServeUnix("root", SocketAddress))
}
