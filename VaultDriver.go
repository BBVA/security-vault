package main

import (
	"bazil.org/fuse"
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"os"
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

	fs, err := NewFS(mountpoint)
	if err != nil {
		fs.errChan <- err
	}

	if err := fs.Mount(r.Name); err != nil {
		fs.errChan <- err
	}

	fmt.Printf("response: %v\n", mountpoint)
	return volume.Response{Mountpoint: mountpoint}

}

func (d VaultDriver) Unmount(r volume.UnmountRequest) volume.Response {

	mountpoint := d.VolumePath + "/" + r.ID + "/" + r.Name

	err := fuse.Unmount(mountpoint)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	fmt.Printf("Unmounted: %s\n", mountpoint)
	return volume.Response{}
}

func (d VaultDriver) Capabilities(r volume.Request) volume.Response {
	return volume.Response{}
}
