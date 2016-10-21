package main

import (
	"fmt"
	"path/filepath"

	"syscall"

	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"descinet.bbva.es/cloudframe-security-vault/utils/fuse"
	"github.com/docker/go-plugins-helpers/volume"
	"golang.org/x/sys/unix"
)

const (
	SocketAddress = "/run/docker/plugins/Vault.sock"
	ServerUrl     = "http://localhost:8200"
	VaultToken    = ""
)

var (
	//DefaultPath = filepath.Join(volume.DefaultDockerRootDirectory, "_vault")
	DefaultMountPath = filepath.Join("/mnt/volumes", "_vault")
	DefaultConfigPath = filepath.Join("/opt", "security-vault")
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
	fuseUtils := fuseutils.NewFuseUtils(fuse)
	dirUtils := filesystem.DefaultDirUtils{}

	driver := NewVaultDriver(DefaultMountPath, ServerUrl, VaultToken, dirUtils, fuseUtils)
	persitor, _ := NewVolumePersistor(DefaultConfigPath, driver, dirUtils)
	handler := volume.NewHandler(persitor)

	fmt.Printf("Listening on %s\n", SocketAddress)
	fmt.Println(handler.ServeUnix("root", SocketAddress))
}

