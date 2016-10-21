package main

import (
	"github.com/docker/go-plugins-helpers/volume"
	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
)

type VolumePersistor struct {
	driver volume.Driver
	dirUtils filesystem.DirUtils
	path   string
}

func NewVolumePersistor(path string, volumeDriver volume.Driver, dirUtils filesystem.DirUtils) (*VolumePersistor, error) {

	_, err := dirUtils.Lstat(path)

	if dirUtils.IsNotExist(err) {
		if err := dirUtils.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	} else if dirUtils.IsExist(err) {
		// TODO recover
	} else if err != nil {
		return nil, err
	}

	return &VolumePersistor{
		driver: volumeDriver,
		dirUtils: dirUtils,
		path: path,
	}, nil
}

func (p VolumePersistor) Create(r volume.Request) volume.Response {
	// TODO save volume data
	return p.driver.Create(r)
}

func (p VolumePersistor) List(r volume.Request) volume.Response {
	return p.driver.List(r)
}

func (p VolumePersistor) Get(r volume.Request) volume.Response {
	return p.driver.Get(r)
}

func (p VolumePersistor) Remove(r volume.Request) volume.Response {
	// TODO remove volume data
	return p.driver.Remove(r)
}

func (p VolumePersistor) Path(r volume.Request) volume.Response {
	return p.driver.Path(r)
}

func (p VolumePersistor) Mount(r volume.MountRequest) volume.Response {
	// TODO save volume data
	return p.driver.Mount(r)
}

func (p VolumePersistor) Unmount(r volume.UnmountRequest) volume.Response {
	// TODO remove volume data
	return p.driver.Unmount(r)
}

func (p VolumePersistor) Capabilities(r volume.Request) volume.Response {
	return p.driver.Capabilities(r)
}
