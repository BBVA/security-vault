package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/docker/go-plugins-helpers/volume"
	"log"
	"os"
	"fmt"
)

type VaultDriver struct {
	VolumePath string
	ServerUrl  string
	VaultToken string
}

func NewVaultDriver(VolumePath string, ServerUrl string, VaultToken string) VaultDriver {
	return VaultDriver{
		VolumePath: VolumePath,
		ServerUrl:  ServerUrl,
		VaultToken: VaultToken,
	}
}

func (d VaultDriver) Create(r volume.Request) volume.Response {

	return volume.Response{}
}

func (d VaultDriver) List(r volume.Request) volume.Response {
	return volume.Response{}
}

func (d VaultDriver) Get(r volume.Request) volume.Response {
	return volume.Response{}
}

func (d VaultDriver) Remove(r volume.Request) volume.Response {
	return volume.Response{}
}

func (d VaultDriver) Path(r volume.Request) volume.Response {
	return volume.Response{}
}

func (d VaultDriver) Mount(r volume.MountRequest) volume.Response {

	mountpoint := d.VolumePath + "/" + r.ID + "/" + r.Name

	fmt.Println("check mountpoint", mountpoint)
	_, err := os.Lstat(mountpoint)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(mountpoint, 0755); err != nil {
			return volume.Response{Err: err.Error()}
		}
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	fmt.Println("mount volume", mountpoint)
	fuseConnection, err := fuse.Mount(
		mountpoint,
		fuse.FSName(r.ID),
		fuse.Subtype("hellofs"),
		fuse.LocalVolume(),
		fuse.VolumeName("Hello world!"),
	)

	if err != nil {
		log.Fatal(err)
	}
	//defer fuseConnection.Close()

	fmt.Println("serve")
	go func() {
		err = fs.Serve(fuseConnection, FS{})
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("expect response from channel Ready")
	// check if the mount process has an error to report
	<-fuseConnection.Ready
	if err := fuseConnection.MountError; err != nil {
		log.Fatal(err)
	}
	fmt.Printf("response: %v\n", mountpoint)
	return volume.Response{Mountpoint: mountpoint}

}

func (d VaultDriver) Unmount(r volume.UnmountRequest) volume.Response {
	return volume.Response{}
}

func (d VaultDriver) Capabilities(r volume.Request) volume.Response {
	return volume.Response{}
}
