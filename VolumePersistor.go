package main

import (
	"github.com/docker/go-plugins-helpers/volume"
	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"path"
	"strings"
	"encoding/json"
	"fmt"
)

type VolumePersistor struct {
	driver    volume.Driver
	dirUtils  filesystem.DirUtils
	fileUtils filesystem.FileUtils
	path      string
}

func NewVolumePersistor(path string, volumeDriver volume.Driver, dirUtils filesystem.DirUtils, fileUtils filesystem.FileUtils) (*VolumePersistor, error) {

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
		fileUtils: fileUtils,
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
	fullFileName := fullFilePath(p.path, r.ID, ".json")

	fileContent, err := mountRequestToJson(r)
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("Could not marshal volume data: %s", err.Error())}
	}
	if err := p.fileUtils.Write(fullFileName, fileContent, 0644); err != nil {
		return volume.Response{Err: fmt.Sprintf("Could not persist volume data: %s", err.Error())}
	}

	return p.driver.Mount(r)
}

func (p VolumePersistor) Unmount(r volume.UnmountRequest) volume.Response {
	fullFileName := fullFilePath(p.path, r.ID, ".json")

	if err := p.dirUtils.RemoveAll(fullFileName); err != nil {
		return volume.Response{Err: fmt.Sprintf("Error removing volume data: %s", err.Error())}
	}

	return p.driver.Unmount(r)
}

func (p VolumePersistor) Capabilities(r volume.Request) volume.Response {
	return p.driver.Capabilities(r)
}

func mountRequestToJson(r volume.MountRequest) ([]byte, error) {
	return json.Marshal(&r)
}

func fullFilePath(dir, name, extension string) string {
	fileName := strings.Join([]string{name, extension}, "")
	return path.Join(dir, fileName)
}