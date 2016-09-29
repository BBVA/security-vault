package main

import (
	"github.com/docker/go-plugins-helpers/volume"
	"os"
	"path"
	"reflect"
	"testing"
	"fmt"
)

type FakeDirUtils struct {
	lstatError    error
	lstatFileInfo os.FileInfo
	mkdirError    error
}

func (f FakeDirUtils) Lstat(mountPoint string) (os.FileInfo, error) {
	return f.lstatFileInfo, f.lstatError
}

func (f FakeDirUtils) MkdirAll(path string, perm os.FileMode) error {
	return f.mkdirError
}

type FakeFuseUtils struct {
	MountError   error
	UnmountError error
}

func (f FakeFuseUtils) Mount(volumeId, mountPoint, volumeName string) error {
	return f.MountError
}

func (f FakeFuseUtils) Unmount(volumeId string) error {
	return f.UnmountError
}

func TestVaultDriver_Mount(t *testing.T) {
	r := volume.MountRequest{
		Name: "Test_volume",
		ID:   "abcdef1234567890",
	}

	fd := FakeDirUtils{
		lstatError:    nil,
		lstatFileInfo: nil,
		mkdirError:    nil,
	}

	ff := FakeFuseUtils{
		MountError:   nil,
		UnmountError: nil,
	}

	d := NewVaultDriver("testpath", "testserver", "testtoken", fd, ff)
	expectedMountPoint := path.Join(d.VolumePath, r.ID, r.Name)

	expectedResponse := volume.Response{Mountpoint: expectedMountPoint}

	response := d.Mount(r)

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Expected %v, received %v\n", expectedResponse, response)
	}
}

func TestVaultDriver_Unmount(t *testing.T) {
	mountRequest := volume.MountRequest{
		Name: "Test_volume",
		ID:   "abcdef1234567890",
	}

	unmountRequest := volume.UnmountRequest{
		Name: "Test_volume",
		ID:   "abcdef1234567890",
	}

	fd := FakeDirUtils{
		lstatError:    nil,
		lstatFileInfo: nil,
		mkdirError:    nil,
	}

	ff := FakeFuseUtils{
		MountError:   nil,
		UnmountError: nil,
	}

	driver := NewVaultDriver("testpath", "testserver", "testtoken", fd, ff)

	driver.Mount(mountRequest)
	response := driver.Unmount(unmountRequest)

	expectedResponse := volume.Response{}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Expected %v, received %v\n", expectedResponse, response)
	}

}
