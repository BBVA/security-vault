package main

import (
	"github.com/docker/go-plugins-helpers/volume"
	"os"
	"path"
	"reflect"
	"testing"
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
	MountError error
}

func (f FakeFuseUtils) Mount(fs *FS, volumeName string) error {
	return f.MountError
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
		MountError: nil,
	}

	d := NewVaultDriver("testpath", "testserver", "testtoken", fd, ff)
	expectedMountPoint := path.Join(d.VolumePath, r.ID, r.Name)

	expectedResponse := volume.Response{Mountpoint: expectedMountPoint}

	response := d.Mount(r)

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Expected %v, received %v\n", expectedResponse, response)
	}
}
