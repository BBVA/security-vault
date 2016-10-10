package main

import (
	"errors"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
	"descinet.bbva.es/cloudframe-security-vault/utils/fuse"
)

type FakeDirUtils struct {
	lstatError    error
	exist         bool
	lstatFileInfo os.FileInfo
	mkdirError    error
}

func (f FakeDirUtils) Lstat(mountPoint string) (os.FileInfo, error) {
	return f.lstatFileInfo, f.lstatError
}

func (f FakeDirUtils) MkdirAll(path string, perm os.FileMode) error {
	return f.mkdirError
}

func (f FakeDirUtils) IsNotExist(err error) bool {
	return !f.exist
}

type Path4FakeFuse struct {
	path string
	err  error
}

type FakeFuseUtils struct {
	MountError   error
	UnmountError error
	CreateError  error
	RemoveError  error
	path         Path4FakeFuse
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

func TestVaultDriver_Mount(t *testing.T) {
	fixtures := []struct {
		mountRequest     volume.MountRequest
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			mountRequest: volume.MountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
			},
			expectedResponse: volume.Response{Mountpoint: path.Join("testpath", "abcdef1234567890", "Test_volume")},
		},
		{
			mountRequest: volume.MountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    errors.New("error"),
				exist:         true,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
			},
			expectedResponse: volume.Response{Err: "error"},
		},
		{
			mountRequest: volume.MountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    errors.New("error"),
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
			},
			expectedResponse: volume.Response{Err: "error"},
		},
		{
			mountRequest: volume.MountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         true,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   errors.New("error"),
				UnmountError: nil,
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		d := NewVaultDriver("testpath", "testserver", "testtoken", fixture.dirUtils, fixture.fuseUtils)
		response := d.Mount(fixture.mountRequest)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_Unmount(t *testing.T) {
	fixtures := []struct {
		unmountRequest   volume.UnmountRequest
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			unmountRequest: volume.UnmountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
			},
			expectedResponse: volume.Response{},
		},
		{
			unmountRequest: volume.UnmountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: errors.New("error"),
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", fixture.dirUtils, fixture.fuseUtils)

		response := driver.Unmount(fixture.unmountRequest)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}

}

func TestVaultDriver_Path(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{
				Name: "Test_volume",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				path: Path4FakeFuse{
					path: "Test_volume",
					err:  nil,
				},
			},
			expectedResponse: volume.Response{Mountpoint: "Test_volume"},
		},
		{
			request: volume.Request{
				Name: "Test_volume",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				path: Path4FakeFuse{
					path: "",
					err:  errors.New("error"),
				},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", fixture.dirUtils, fixture.fuseUtils)

		response := driver.Path(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_Capabilities(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				path: Path4FakeFuse{
					path: "Test_volume",
					err:  nil,
				},
			},
			expectedResponse: volume.Response{
				Capabilities: volume.Capability{
					Scope: "local",
				},
			},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", fixture.dirUtils, fixture.fuseUtils)

		response := driver.Capabilities(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_Create(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{
				Name: "test_vol",
				Options: map[string]string{
					"key1": "val1",
				},
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				CreateError: nil,
				path: Path4FakeFuse{
					path: "Test_volume",
					err:  nil,
				},
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{
				Name: "test_vol",
				Options: map[string]string{
					"key1": "val1",
				},
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				CreateError: errors.New("error"),
				path: Path4FakeFuse{
					path: "Test_volume",
					err:  nil,
				},
			},
			expectedResponse: volume.Response{Err:"error"},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", fixture.dirUtils, fixture.fuseUtils)

		response := driver.Create(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_Remove(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{
				Name: "test_vol",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				RemoveError: nil,
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{
				Name: "test_vol",
				Options: map[string]string{
					"key1": "val1",
				},
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				RemoveError: errors.New("error"),
							},
			expectedResponse: volume.Response{Err:"error"},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", fixture.dirUtils, fixture.fuseUtils)

		response := driver.Remove(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}
