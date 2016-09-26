package main

import (
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"os"
	"path"
)

type DirUtils interface {
	Lstat(mountPoint string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
}

type DefaultDirUtils struct {}

func (d DefaultDirUtils) MkdirAll(path string, perm os.FileMode) (error) {
	return os.MkdirAll(path,perm)
}

func (d DefaultDirUtils) Lstat(mountPoint string) (os.FileInfo, error) {
	return os.Lstat(mountPoint)
}

type VaultDriver struct {
	VolumePath string
	ServerUrl  string
	VaultToken string
	fs         map[string]*FS
	dirUtils   DirUtils
}

func NewVaultDriver(VolumePath string, ServerUrl string, VaultToken string, dirUtils DirUtils) VaultDriver {
	return VaultDriver{
		VolumePath: VolumePath,
		ServerUrl:  ServerUrl,
		VaultToken: VaultToken,
		fs: map[string]*FS{},
		dirUtils: &dirUtils,
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

	fs, err := NewFS(mountPoint)
	if err != nil {
		fs.errChan <- err
	}

	if err := fs.Mount(r.Name); err != nil {
		fs.errChan <- err
	}

	d.fs[r.ID] = fs

	fmt.Printf("response: %v\n", mountPoint)
	return volume.Response{Mountpoint: mountPoint}

}

func (d VaultDriver) Unmount(r volume.UnmountRequest) volume.Response {
	if fs, ok := d.fs[r.ID]; ok {
		err := fs.Unmount()
		if err != nil {
			return volume.Response{Err: err.Error()}
		}

		fmt.Printf("Unmounted: %s\n", fs.mountpoint)
		return volume.Response{}
	} else {
		return volume.Response{Err: "Volume not found"}
	}

}

func (d VaultDriver) Capabilities(r volume.Request) volume.Response {
	return volume.Response{}
}
