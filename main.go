package main

import (
	"fmt"
	"path/filepath"

	"descinet.bbva.es/cloudframe-security-vault/utils/fuse"

	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"github.com/docker/go-plugins-helpers/volume"
)

const (
	SocketAddress = "/run/docker/plugins/Vault.sock"
	ServerUrl     = "http://localhost:8200"
	VaultToken    = ""
)

var (
	//DefaultPath = filepath.Join(volume.DefaultDockerRootDirectory, "_vault")
	DefaultPath = filepath.Join("/mnt/volumes", "_vault")
)

/* hay que implementar argumentos para recibir:
*
*  -path de la configuración ( aquí pondremos todo lo que ahora está como constantes )
*
 */

func main() {

	dirUtils := filesystem.DefaultDirUtils{}
	fuseUtils := fuseutils.NewFuseUtils()
	d := NewVaultDriver(DefaultPath, ServerUrl, VaultToken, dirUtils, fuseUtils)
	h := volume.NewHandler(d)
	fmt.Printf("Listening on %s\n", SocketAddress)
	fmt.Println(h.ServeUnix("root", SocketAddress))

}
