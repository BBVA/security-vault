package main

import (
	"github.com/docker/go-plugins-helpers/volume"
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
	return volume.Response{}
}

func (d VaultDriver) Unmount(r volume.UnmountRequest) volume.Response {
	return volume.Response{}
}

func (d VaultDriver) Capabilities(r volume.Request) volume.Response {
	return volume.Response{}
}
