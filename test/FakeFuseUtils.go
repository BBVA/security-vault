package test

import "descinet.bbva.es/cloudframe-security-vault/utils/fuse"


type Path4FakeFuse struct {
	path string
	err  error
}

type List4FakeFuse struct {
	volumes []fuseutils.VolumeData
	err     error
}

type Get4FakeFuse struct {
	volume fuseutils.VolumeData
	err    error
}

type FakeFuseUtils struct {
	MountError   error
	UnmountError error
	CreateError  error
	RemoveError  error
	list         List4FakeFuse
	path         Path4FakeFuse
	get          Get4FakeFuse
}

func (f FakeFuseUtils) Mount(volumeId, mountPoint, volumeName string) error {
	return f.MountError
}

func (f FakeFuseUtils) Unmount(volumeName string) error {
	return f.UnmountError
}

func (f FakeFuseUtils) Path(volumeName string) (string, error) {
	return f.path.path, f.path.err
}

func (f FakeFuseUtils) Create(volumeName string, options fuseutils.VolumeOptions) error {
	return f.CreateError
}

func (f FakeFuseUtils) Remove(volumeName string) error {
	return f.RemoveError
}

func (f FakeFuseUtils) List() ([]fuseutils.VolumeData, error) {
	return f.list.volumes, f.list.err
}

func (f FakeFuseUtils) Get(volumeName string) (fuseutils.VolumeData, error) {
	return f.get.volume, f.get.err
}
