package main

import (
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"os"
	"path"
)

type VaultDriver struct {
	VolumePath string
	ServerUrl  string
	VaultToken string
	dirUtils   DirUtils
	fuseUtils  FuseUtils
}

func NewVaultDriver(VolumePath string, ServerUrl string, VaultToken string, dirUtils DirUtils, fuseUtils FuseUtils) VaultDriver {
	fmt.Println(fuseUtils)
	return VaultDriver{
		VolumePath: VolumePath,
		ServerUrl:  ServerUrl,
		VaultToken: VaultToken,
		dirUtils:   dirUtils,
		fuseUtils:  fuseUtils,
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

	mountPoint := path.Join(d.VolumePath, r.ID, r.Name)

	fmt.Println("check mountpoint", mountPoint)
	_, err := d.dirUtils.Lstat(mountPoint)

	if os.IsNotExist(err) {
		if err := d.dirUtils.MkdirAll(mountPoint, 0755); err != nil {
			return volume.Response{Err: err.Error()}
		}
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	fmt.Println("mount volume", mountPoint)

	if err := d.fuseUtils.Mount(r.ID, mountPoint, r.Name); err != nil {
		//fs.errChan <- err
	}

	fmt.Printf("response: %v\n", mountPoint)
	return volume.Response{Mountpoint: mountPoint}

}

func (d VaultDriver) Unmount(r volume.UnmountRequest) volume.Response {

	err := d.fuseUtils.Unmount(r.ID)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	fmt.Printf("Unmounted: %s\n", r.ID)
	return volume.Response{}

}

func (d VaultDriver) Capabilities(r volume.Request) volume.Response {
	return volume.Response{}
}
