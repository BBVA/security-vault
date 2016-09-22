package main

import (
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"path/filepath"
)

const (
	SocketAddress = "/run/docker/plugins/Vault.sock"
	ServerUrl     = "http://localhost:8200"
	VaultToken    = ""
)

var (
	DefaultPath = filepath.Join(volume.DefaultDockerRootDirectory, "_vault")
)

/* hay que implementar argumentos para recibir:
*
*  -path de la configuración ( aquí pondremos todo lo que ahora está como constantes )
*
 */

func main() {

	d := NewVaultDriver(DefaultPath, ServerUrl, VaultToken)
	h := volume.NewHandler(d)
	fmt.Printf("Listening on %s\n", SocketAddress)
	fmt.Println(h.ServeUnix("root", SocketAddress))

}
