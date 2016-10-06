package main

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
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
	fs           map[string]*filesystem.FS
}

func (f FakeFuseUtils) Mount(volumeId, mountPoint, volumeName string) error {
	f.fs[volumeName].VolumeId = volumeId
	f.fs[volumeName].Mountpoint = mountPoint
	return f.MountError
}

func (f FakeFuseUtils) Unmount(volumeName string) error {
	return f.UnmountError
}

func (f FakeFuseUtils) Path(volumeName string) (string, error) {
	return f.fs[volumeName].Mountpoint, nil
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

	fs := make(map[string]*filesystem.FS)

	fs[r.Name] = &filesystem.FS{}

	ff := FakeFuseUtils{
		MountError:   nil,
		UnmountError: nil,
		fs:           fs,
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

	fs := make(map[string]*filesystem.FS)

	fs[mountRequest.Name] = &filesystem.FS{}

	ff := FakeFuseUtils{
		MountError:   nil,
		UnmountError: nil,
		fs:           fs,
	}

	driver := NewVaultDriver("testpath", "testserver", "testtoken", fd, ff)

	driver.Mount(mountRequest)
	response := driver.Unmount(unmountRequest)

	expectedResponse := volume.Response{}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Expected %v, received %v\n", expectedResponse, response)
	}

}

func TestVaultDriver_Path(t *testing.T) {

	mountRequest := volume.MountRequest{
		Name: "Test_volume",
		ID:   "abcdef1234567890",
	}

	request := volume.Request{
		Name: "Test_volume",
	}

	fd := FakeDirUtils{
		lstatError:    nil,
		lstatFileInfo: nil,
		mkdirError:    nil,
	}

	fs := make(map[string]*filesystem.FS)

	fs[mountRequest.Name] = &filesystem.FS{}

	ff := FakeFuseUtils{
		MountError:   nil,
		UnmountError: nil,
		fs:           fs,
	}

	driver := NewVaultDriver("testpath", "testserver", "testtoken", fd, ff)

	mountresponse := driver.Mount(mountRequest)
	response := driver.Path(request)

	if !reflect.DeepEqual(response.Mountpoint, mountresponse.Mountpoint) {
		t.Errorf("Expected %v, received %v\n", mountresponse.Mountpoint, response.Mountpoint)
	}
}
